package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/api/session"
	"github.com/example/go-svc-boilerplate/internal/cnst"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/core/authorize"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// RequireAuth guards routes with the sealed session cookie. It unseals the
// cookie, validates (and, if expired, refreshes) the WorkOS access token, sets
// a refreshed cookie when needed, and stashes the current user under
// cnst.CtxUserKey. Any failure clears the cookie and aborts with 401.
func RequireAuth(auth *core.Auth, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens, err := session.Read(c, cfg)
		if err != nil {
			session.Clear(c, cfg)
			abort(c, errs.NewUnauthorized("no session cookie: "+err.Error(), "not authenticated"))
			return
		}

		azCtx := &authorize.AuthorizeCtx{
			Ctx:          auth.BaseCtx(cnst.LangEN),
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		if err := authorize.New(azCtx).Do(&core.DoCtx{}); err != nil {
			session.Clear(c, cfg)
			abort(c, err)
			return
		}

		// Access token was refreshed — write the updated session cookie.
		if azCtx.NewTokens != nil {
			if err := session.Set(c, cfg, azCtx.NewTokens); err != nil {
				abort(c, errs.NewInternalServer("failed to update session cookie: "+err.Error(), ""))
				return
			}
		}

		c.Set(cnst.CtxUserKey, azCtx.User)
		c.Next()
	}
}

// abort renders an errs error and stops the handler chain.
func abort(c *gin.Context, e error) {
	var coder errs.StatusCoder
	var formatter errs.HTTPFormatter
	if ce, ok := e.(*errs.Error); ok {
		coder, formatter = ce, ce
		c.AbortWithStatusJSON(coder.StatusCode(), formatter.HTTPFormat())
		return
	}
	c.AbortWithStatusJSON(500, gin.H{"status": false, "error": "something went wrong!"})
}
