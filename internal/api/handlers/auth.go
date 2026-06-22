package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/api/session"
	"github.com/example/go-svc-boilerplate/internal/cnst"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/core/logouturl"
	"github.com/example/go-svc-boilerplate/internal/core/signupcb"
	"github.com/example/go-svc-boilerplate/internal/core/signupurl"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// AuthHandler is the HTTP entry point for the WorkOS AuthKit signup flows. It
// binds the request, runs the flow, and serializes the result — no business
// logic lives here.
type AuthHandler struct {
	auth *core.Auth
	cfg  *config.Config
}

func NewAuthHandler(auth *core.Auth, cfg *config.Config) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg}
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
	if err := session.Set(c, h.cfg, ctx.Tokens); err != nil {
		ServeErr(c, err)
		return
	}
	ServeData(c, ctx.Resp)
}

// Logout reads the session cookie, returns the AuthKit hosted logout URL for
// ending the WorkOS session, and clears the local session cookie.
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := &logouturl.LogoutURLCtx{
		Ctx:      h.auth.BaseCtx(lang(c)),
		ReturnTo: c.Query("return_to"),
	}
	if tokens, err := session.Read(c, h.cfg); err == nil {
		ctx.AccessToken = tokens.AccessToken
	}
	if err := logouturl.New(ctx).Do(&core.DoCtx{}); err != nil {
		ServeErr(c, err)
		return
	}
	session.Clear(c, h.cfg)
	ServeData(c, ctx.Resp)
}

// Me returns the currently authenticated user. RequireAuth must run first; it
// stashes the user under cnst.CtxUserKey.
func (h *AuthHandler) Me(c *gin.Context) {
	user, ok := c.Get(cnst.CtxUserKey)
	if !ok {
		ServeErr(c, errs.NewUnauthorized("no user in context", "not authenticated"))
		return
	}
	u, ok := user.(*entity.User)
	if !ok {
		ServeErr(c, errs.NewUnauthorized("invalid user in context", "not authenticated"))
		return
	}
	ServeData(c, h.auth.TransformUser(u))
}
