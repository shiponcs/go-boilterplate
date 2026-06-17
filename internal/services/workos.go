package services

import (
	"context"
	"fmt"

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
	log         *zap.Logger
}

func NewWorkOSAuthService(cfg *config.Config, log *zap.Logger) *workosAuthService {
	usermanagement.SetAPIKey(cfg.WorkOS.ApiKey)
	return &workosAuthService{
		clientID:    cfg.WorkOS.ClientID,
		redirectURI: cfg.WorkOS.RedirectURI,
		log:         log,
	}
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
	}, nil
}
