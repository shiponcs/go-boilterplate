package signupcb

import (
	"context"
	"strings"

	"github.com/example/go-svc-boilerplate/internal/cnst"
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/localization"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// validate rejects callbacks missing a code or carrying a state we did not
// issue (or that was already used). ConsumeSignupState is one-shot.
type validate struct {
	ctx *SignupCbCtx
}

func (v *validate) Do(*core.DoCtx) error {
	if strings.TrimSpace(v.ctx.Code) == "" {
		msg := localization.GetMessage("auth_invalid_code", v.ctx.Lang)
		return errs.NewBadReq("missing authorization code", msg)
	}

	ok, err := v.ctx.Cache.ConsumeSignupState(v.ctx.State)
	if err != nil {
		return errs.NewInternalServer("failed to verify state: "+err.Error(), "")
	}
	if !ok {
		msg := localization.GetMessage("auth_state_mismatch", v.ctx.Lang)
		return errs.NewBadReq("invalid or expired state", msg)
	}
	return nil
}

// exchange trades the authorization code for the authenticated WorkOS user.
type exchange struct {
	ctx *SignupCbCtx
}

func (e *exchange) Do(*core.DoCtx) error {
	result, err := e.ctx.Srv.AuthService.AuthenticateWithCode(context.Background(), e.ctx.Code)
	if err != nil {
		msg := localization.GetMessage("signup_failed", e.ctx.Lang)
		return errs.NewInternalServer("failed to authenticate with code: "+err.Error(), msg)
	}
	e.ctx.result = result
	return nil
}

// persist upserts the local user mirror keyed by workos_user_id.
type persist struct {
	ctx *SignupCbCtx
}

func (p *persist) Do(*core.DoCtx) error {
	r := p.ctx.result
	user := &entity.User{
		WorkOSUserID:  r.WorkOSUserID,
		Email:         r.Email,
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		EmailVerified: r.EmailVerified,
		Status:        cnst.UserStatusActive,
	}
	if err := p.ctx.Store.UserStore.Upsert(user); err != nil {
		return errs.NewInternalServer("failed to persist user: "+err.Error(), "")
	}
	p.ctx.User = user
	return nil
}

// response builds the API payload from the persisted user.
type response struct {
	ctx *SignupCbCtx
}

func (r *response) Do(*core.DoCtx) error {
	u := r.ctx.User
	r.ctx.Resp = &dto.UserResponse{
		ID:            u.ID,
		WorkOSUserID:  u.WorkOSUserID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		EmailVerified: u.EmailVerified,
		Status:        u.Status,
		CreatedAt:     u.CreatedAt.Unix(),
	}
	return nil
}
