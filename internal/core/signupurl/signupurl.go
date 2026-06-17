package signupurl

import "github.com/example/go-svc-boilerplate/internal/core"

// signupURL assembles the ordered steps for issuing an AuthKit sign-up URL:
//
//	gen state -> build url -> build response
type signupURL struct {
	ctx *SignupURLCtx
}

func New(ctx *SignupURLCtx) core.Doer {
	return &signupURL{ctx: ctx}
}

func (s *signupURL) Do(doCtx *core.DoCtx) error {
	doers := core.Doers{
		&genState{ctx: s.ctx},
		&buildURL{ctx: s.ctx},
		&response{ctx: s.ctx},
	}
	return doers.Do(doCtx)
}
