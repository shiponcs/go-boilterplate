package core

// Doer is one small step in a business flow.
type Doer interface {
	Do(*DoCtx) error
}

// DoCtx carries pipeline control flags shared across steps.
type DoCtx struct {
	// IsExit short-circuits the rest of the pipeline (the loop breaks).
	IsExit bool
	// NxtDoer jumps execution to / resumes at a specific step; steps before
	// it are skipped. Useful for branching or resumable flows.
	NxtDoer Doer
}

// Doers is an ordered list of steps with its own Do method, so a flow is just a
// slice of small Doers assembled and run in order.
type Doers []Doer

func (ds Doers) Do(ctx *DoCtx) error {
	for _, d := range ds {
		if ctx != nil && ctx.NxtDoer != nil && ctx.NxtDoer != d {
			continue
		}
		ctx.NxtDoer = nil
		if err := d.Do(ctx); err != nil {
			return err
		}
		if ctx != nil && ctx.IsExit {
			break
		}
	}
	return nil
}
