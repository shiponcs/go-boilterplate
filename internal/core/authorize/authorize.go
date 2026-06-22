package authorize

import "github.com/example/go-svc-boilerplate/internal/core"

// authorize assembles the ordered steps for resolving the current session:
//
//	resolve claims -> refresh if expired -> load user
type authorize struct {
	ctx *AuthorizeCtx
}

func New(ctx *AuthorizeCtx) core.Doer {
	return &authorize{ctx: ctx}
}

func (a *authorize) Do(doCtx *core.DoCtx) error {
	doers := core.Doers{
		&resolveClaims{ctx: a.ctx},
		&refreshIfExpired{ctx: a.ctx},
		&loadUser{ctx: a.ctx},
	}
	return doers.Do(doCtx)
}
