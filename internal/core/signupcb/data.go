package signupcb

import (
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/models/value"
)

// SignupCbCtx embeds the shared core.Ctx and adds fields specific to handling
// the AuthKit callback (code exchange + local user upsert).
type SignupCbCtx struct {
	core.Ctx

	Code  string
	State string

	result *value.AuthnResult
	User   *entity.User
	Resp   *dto.UserResponse
	// Tokens carries the WorkOS session tokens to the handler, which seals them
	// into the session cookie.
	Tokens *value.SessionTokens
}
