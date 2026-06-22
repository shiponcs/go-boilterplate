package authorize

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// resolveClaims verifies the access token. A valid token populates claims; an
// expired one leaves claims nil so refreshIfExpired runs; anything else is a
// hard 401.
type resolveClaims struct {
	ctx *AuthorizeCtx
}

func (s *resolveClaims) Do(*core.DoCtx) error {
	if s.ctx.AccessToken == "" {
		return errs.NewUnauthorized("missing access token", "not authenticated")
	}
	claims, err := s.ctx.Srv.AuthService.ValidateAccessToken(s.ctx.AccessToken)
	if err == nil {
		s.ctx.claims = claims
		return nil
	}
	if errors.Is(err, services.ErrTokenExpired) {
		return nil // claims stay nil; refreshIfExpired handles it
	}
	return errs.NewUnauthorized("invalid access token: "+err.Error(), "not authenticated")
}

// refreshIfExpired exchanges the refresh token for a new token pair when the
// access token was expired, and re-derives the claims from the fresh token.
type refreshIfExpired struct {
	ctx *AuthorizeCtx
}

func (s *refreshIfExpired) Do(*core.DoCtx) error {
	if s.ctx.claims != nil {
		return nil // access token was still valid
	}
	if s.ctx.RefreshToken == "" {
		return errs.NewUnauthorized("expired token, no refresh token", "session expired")
	}
	tokens, err := s.ctx.Srv.AuthService.RefreshSession(context.Background(), s.ctx.RefreshToken)
	if err != nil {
		return errs.NewUnauthorized("refresh failed: "+err.Error(), "session expired")
	}
	claims, err := s.ctx.Srv.AuthService.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		return errs.NewUnauthorized("refreshed token invalid: "+err.Error(), "session expired")
	}
	s.ctx.NewTokens = tokens
	s.ctx.claims = claims
	return nil
}

// loadUser fetches the local user mirror identified by the token's sub claim.
type loadUser struct {
	ctx *AuthorizeCtx
}

func (s *loadUser) Do(*core.DoCtx) error {
	user, err := s.ctx.Store.UserStore.GetByWorkOSID(s.ctx.claims.WorkOSUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewUnauthorized("user not found for token", "not authenticated")
		}
		return errs.NewInternalServer("failed to load user: "+err.Error(), "")
	}
	s.ctx.User = user
	return nil
}
