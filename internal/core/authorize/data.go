package authorize

import (
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/models/value"
)

// AuthorizeCtx embeds the shared core.Ctx and carries the session tokens read
// from the cookie. It resolves the current user, refreshing the access token
// when it has expired.
type AuthorizeCtx struct {
	core.Ctx

	AccessToken  string
	RefreshToken string

	claims *value.SessionClaims
	// NewTokens is non-nil when the access token was refreshed; the caller
	// re-seals it into the session cookie.
	NewTokens *value.SessionTokens
	User      *entity.User
}
