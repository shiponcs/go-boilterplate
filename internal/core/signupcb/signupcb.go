package signupcb

import "github.com/example/go-svc-boilerplate/internal/core"

// signupCb assembles the ordered steps for the AuthKit callback:
//
//	validate -> exchange code -> persist user -> build response
type signupCb struct {
	ctx *SignupCbCtx
}

func New(ctx *SignupCbCtx) core.Doer {
	return &signupCb{ctx: ctx}
}

func (s *signupCb) Do(doCtx *core.DoCtx) error {
	doers := core.Doers{
		&validate{ctx: s.ctx},
		&exchange{ctx: s.ctx},
		&persist{ctx: s.ctx},
		&response{ctx: s.ctx},
	}
	return doers.Do(doCtx)
}
