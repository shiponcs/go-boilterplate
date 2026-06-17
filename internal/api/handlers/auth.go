package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/core/signupcb"
	"github.com/example/go-svc-boilerplate/internal/core/signupurl"
)

// AuthHandler is the HTTP entry point for the WorkOS AuthKit signup flows. It
// binds the request, runs the flow, and serializes the result — no business
// logic lives here.
type AuthHandler struct {
	auth *core.Auth
}

func NewAuthHandler(auth *core.Auth) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Signup returns the AuthKit hosted authorization URL the client should open to
// complete sign-up.
func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := &signupurl.SignupURLCtx{
		Ctx: h.auth.BaseCtx(lang(c)),
	}
	if err := signupurl.New(ctx).Do(&core.DoCtx{}); err != nil {
		ServeErr(c, err)
		return
	}
	ServeData(c, ctx.Resp)
}

// Callback handles the WorkOS redirect: it validates state, exchanges the
// authorization code, mirrors the user locally, and returns it.
func (h *AuthHandler) Callback(c *gin.Context) {
	ctx := &signupcb.SignupCbCtx{
		Ctx:   h.auth.BaseCtx(lang(c)),
		Code:  c.Query("code"),
		State: c.Query("state"),
	}
	if err := signupcb.New(ctx).Do(&core.DoCtx{}); err != nil {
		ServeErr(c, err)
		return
	}
	ServeData(c, ctx.Resp)
}
