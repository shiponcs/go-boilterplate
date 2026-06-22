package services

import (
	"context"

	"github.com/example/go-svc-boilerplate/internal/models/value"
)

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
	// authenticated user.
	AuthenticateWithCode(ctx context.Context, code string) (*value.AuthnResult, error)
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
