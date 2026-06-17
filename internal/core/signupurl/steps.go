package signupurl

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// stateTTL bounds how long an issued sign-up URL stays valid before the user
// must request a fresh one.
const stateTTL = 10 * time.Minute

// genState creates a random CSRF state and records it in the cache so the
// callback can verify it was issued by us.
type genState struct {
	ctx *SignupURLCtx
}

func (g *genState) Do(*core.DoCtx) error {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return errs.NewInternalServer("failed to generate state: "+err.Error(), "")
	}
	state := hex.EncodeToString(buf)
	if err := g.ctx.Cache.SetSignupState(state, stateTTL); err != nil {
		return errs.NewInternalServer("failed to store state: "+err.Error(), "")
	}
	g.ctx.State = state
	return nil
}

// buildURL asks the auth service for the AuthKit sign-up authorization URL.
type buildURL struct {
	ctx *SignupURLCtx
}

func (b *buildURL) Do(*core.DoCtx) error {
	url, err := b.ctx.Srv.AuthService.SignupURL(b.ctx.State)
	if err != nil {
		return errs.NewInternalServer("failed to build signup url: "+err.Error(), "")
	}
	b.ctx.url = url
	return nil
}

// response builds the API payload from the generated URL.
type response struct {
	ctx *SignupURLCtx
}

func (r *response) Do(*core.DoCtx) error {
	r.ctx.Resp = &dto.SignupURLResponse{AuthorizationURL: r.ctx.url}
	return nil
}
