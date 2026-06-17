package signupurl

import (
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
)

// SignupURLCtx embeds the shared core.Ctx and adds fields for building the
// AuthKit sign-up authorization URL.
type SignupURLCtx struct {
	core.Ctx

	State string
	url   string
	Resp  *dto.SignupURLResponse
}
