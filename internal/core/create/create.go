package create

import "github.com/example/go-svc-boilerplate/internal/core"

// create assembles the ordered steps for creating a widget. The flow shape is
// the convention used across the service:
//
//	validate -> compute -> persist -> build response
//
// Add side-effect steps (publish event, enqueue job, start workflow) after
// persist as needed.
type create struct {
	ctx *CreateCtx
}

func New(ctx *CreateCtx) core.Doer {
	return &create{ctx: ctx}
}

func (c *create) Do(doCtx *core.DoCtx) error {
	doers := core.Doers{
		&validate{ctx: c.ctx},
		&quote{ctx: c.ctx},
		&persist{ctx: c.ctx},
		&response{ctx: c.ctx},
	}
	return doers.Do(doCtx)
}
