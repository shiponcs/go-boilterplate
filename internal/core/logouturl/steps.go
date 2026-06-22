package logouturl

import (
	"strings"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// validate requires a session id because WorkOS logout revokes a concrete
// session identified by sid.
type validate struct {
	ctx *LogoutURLCtx
}

func (v *validate) Do(*core.DoCtx) error {
	if strings.TrimSpace(v.ctx.SessionID) == "" {
		return errs.NewBadReq("missing session_id", "session_id is required")
	}
	return nil
}

// buildURL asks the auth service for the AuthKit logout URL.
type buildURL struct {
	ctx *LogoutURLCtx
}

func (b *buildURL) Do(*core.DoCtx) error {
	url, err := b.ctx.Srv.AuthService.LogoutURL(b.ctx.SessionID, b.ctx.ReturnTo)
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
