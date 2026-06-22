package logouturl

import (
	"strings"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// validate requires an access token and extracts the WorkOS session id (sid)
// from it, since WorkOS logout revokes a concrete session identified by sid.
type validate struct {
	ctx *LogoutURLCtx
}

func (v *validate) Do(*core.DoCtx) error {
	if strings.TrimSpace(v.ctx.AccessToken) == "" {
		return errs.NewUnauthorized("missing session", "not authenticated")
	}
	sid, err := v.ctx.Srv.AuthService.SessionIDFromToken(v.ctx.AccessToken)
	if err != nil {
		return errs.NewUnauthorized("cannot read session id: "+err.Error(), "not authenticated")
	}
	v.ctx.sid = sid
	return nil
}

// buildURL asks the auth service for the AuthKit logout URL.
type buildURL struct {
	ctx *LogoutURLCtx
}

func (b *buildURL) Do(*core.DoCtx) error {
	url, err := b.ctx.Srv.AuthService.LogoutURL(b.ctx.sid, b.ctx.ReturnTo)
	if err != nil {
		return errs.NewInternalServer("failed to build logout url: "+err.Error(), "")
	}
	b.ctx.url = url
	return nil
}

// response builds the API payload from the generated URL.
type response struct {
	ctx *LogoutURLCtx
}

func (r *response) Do(*core.DoCtx) error {
	r.ctx.Resp = &dto.LogoutURLResponse{LogoutURL: r.ctx.url}
	return nil
}
