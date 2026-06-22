package services

import (
	"context"
	"errors"

	"github.com/example/go-svc-boilerplate/internal/models/value"
)

// ErrTokenExpired signals that an access token failed validation solely
// because it has expired, so the caller can attempt a refresh.
var ErrTokenExpired = errors.New("access token expired")

// PricingService is the interface core code depends on. Defining the interface
// here (consumer side) and providing a concrete impl via fx lets you swap or
// mock the client without touching callers.
type PricingService interface {
	Quote(units int) (*value.CalculatedPrice, error)
}

// AuthService is the identity provider surface core code depends on. The
// concrete impl wraps the WorkOS AuthKit SDK.
type AuthService interface {
	// SignupURL builds the AuthKit hosted authorization URL for sign-up,
	// carrying the given CSRF state.
	SignupURL(state string) (string, error)
	// LogoutURL builds the AuthKit logout URL for ending a session.
	LogoutURL(sessionID, returnTo string) (string, error)
	// AuthenticateWithCode exchanges an AuthKit authorization code for the
	// authenticated user (including the session tokens).
	AuthenticateWithCode(ctx context.Context, code string) (*value.AuthnResult, error)
	// ValidateAccessToken verifies a WorkOS access-token JWT against the JWKS
	// (signature + expiry) and returns its claims. It returns ErrTokenExpired
	// when the token is well-formed but expired (the caller should refresh).
	ValidateAccessToken(token string) (*value.SessionClaims, error)
	// RefreshSession exchanges a refresh token for a new token pair.
	RefreshSession(ctx context.Context, refreshToken string) (*value.SessionTokens, error)
	// SessionIDFromToken extracts the WorkOS session id (sid) from an access
	// token without verifying it — used to build the logout URL.
	SessionIDFromToken(token string) (string, error)
}

// SrvHolder aggregates every external service client into one injectable
// struct, mirroring StoHolder. Add a field per new service.
type SrvHolder struct {
	PricingService PricingService
	AuthService    AuthService
}

func NewSrvHolder(pricing PricingService, auth AuthService) *SrvHolder {
	return &SrvHolder{
		PricingService: pricing,
		AuthService:    auth,
	}
}
