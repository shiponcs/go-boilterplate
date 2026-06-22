package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/value"
)

// workosAuthService implements AuthService against WorkOS AuthKit. The SDK uses
// a package-level default client configured via SetAPIKey, done once here.
type workosAuthService struct {
	clientID    string
	redirectURI string
	jwks        keyfunc.Keyfunc // caches/rotates the WorkOS JWKS for local validation
	log         *zap.Logger
}

func NewWorkOSAuthService(cfg *config.Config, log *zap.Logger) *workosAuthService {
	usermanagement.SetAPIKey(cfg.WorkOS.ApiKey)

	s := &workosAuthService{
		clientID:    cfg.WorkOS.ClientID,
		redirectURI: cfg.WorkOS.RedirectURI,
		log:         log,
	}

	// Build the JWKS keyfunc used to verify access tokens locally. Best-effort:
	// if it can't be built (e.g. no client id / offline), log and leave it nil —
	// ValidateAccessToken then fails closed rather than blocking startup.
	if cfg.WorkOS.ClientID != "" {
		jwksURL, err := usermanagement.GetJWKSURL(cfg.WorkOS.ClientID)
		if err != nil {
			log.Warn("could not build WorkOS JWKS url", zap.Error(err))
		} else if k, err := keyfunc.NewDefault([]string{jwksURL.String()}); err != nil {
			log.Warn("could not initialize WorkOS JWKS", zap.Error(err))
		} else {
			s.jwks = k
		}
	}

	return s
}

func (s *workosAuthService) SignupURL(state string) (string, error) {
	u, err := usermanagement.GetAuthorizationURL(usermanagement.GetAuthorizationURLOpts{
		ClientID:    s.clientID,
		RedirectURI: s.redirectURI,
		Provider:    "authkit",
		State:       state,
		ScreenHint:  usermanagement.SignUp,
	})
	if err != nil {
		return "", fmt.Errorf("build authkit url: %w", err)
	}
	return u.String(), nil
}

func (s *workosAuthService) LogoutURL(sessionID, returnTo string) (string, error) {
	u, err := usermanagement.GetLogoutURL(usermanagement.GetLogoutURLOpts{
		SessionID: sessionID,
		ReturnTo:  returnTo,
	})
	if err != nil {
		return "", fmt.Errorf("build logout url: %w", err)
	}
	return u.String(), nil
}

func (s *workosAuthService) AuthenticateWithCode(ctx context.Context, code string) (*value.AuthnResult, error) {
	resp, err := usermanagement.AuthenticateWithCode(ctx, usermanagement.AuthenticateWithCodeOpts{
		ClientID: s.clientID,
		Code:     code,
	})
	if err != nil {
		return nil, fmt.Errorf("authenticate with code: %w", err)
	}
	return &value.AuthnResult{
		WorkOSUserID:  resp.User.ID,
		Email:         resp.User.Email,
		FirstName:     resp.User.FirstName,
		LastName:      resp.User.LastName,
		EmailVerified: resp.User.EmailVerified,
		AccessToken:   resp.AccessToken,
		RefreshToken:  resp.RefreshToken,
	}, nil
}

// ValidateAccessToken verifies the JWT against the WorkOS JWKS and returns its
// claims. An expired-but-otherwise-valid token yields ErrTokenExpired so the
// caller can refresh.
func (s *workosAuthService) ValidateAccessToken(token string) (*value.SessionClaims, error) {
	if s.jwks == nil {
		return nil, errors.New("jwks not initialized")
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, s.jwks.Keyfunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("validate access token: %w", err)
	}
	return claimsToSession(claims), nil
}

// RefreshSession exchanges a refresh token for a fresh token pair.
func (s *workosAuthService) RefreshSession(ctx context.Context, refreshToken string) (*value.SessionTokens, error) {
	resp, err := usermanagement.AuthenticateWithRefreshToken(ctx, usermanagement.AuthenticateWithRefreshTokenOpts{
		ClientID:     s.clientID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh session: %w", err)
	}
	return &value.SessionTokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// SessionIDFromToken pulls the sid claim out of an access token without
// verifying it (the token came from a sealed cookie we issued).
func (s *workosAuthService) SessionIDFromToken(token string) (string, error) {
	claims := jwt.MapClaims{}
	if _, _, err := jwt.NewParser().ParseUnverified(token, claims); err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	sid, _ := claims["sid"].(string)
	if sid == "" {
		return "", errors.New("token has no sid claim")
	}
	return sid, nil
}

// claimsToSession maps WorkOS JWT claims to the domain value object.
func claimsToSession(claims jwt.MapClaims) *value.SessionClaims {
	sub, _ := claims["sub"].(string)
	sid, _ := claims["sid"].(string)
	var exp int64
	if v, err := claims.GetExpirationTime(); err == nil && v != nil {
		exp = v.Unix()
	}
	return &value.SessionClaims{
		WorkOSUserID: sub,
		WorkOSSID:    sid,
		Expiry:       exp,
	}
}
