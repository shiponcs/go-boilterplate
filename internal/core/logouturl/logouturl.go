package logouturl

import "github.com/example/go-svc-boilerplate/internal/core"

// logoutURL assembles the ordered steps for issuing an AuthKit logout URL:
//
//	validate -> build url -> build response
type logoutURL struct {
	ctx *LogoutURLCtx
}

func New(ctx *LogoutURLCtx) core.Doer {
	return &logoutURL{ctx: ctx}
}

func (s *logoutURL) Do(doCtx *core.DoCtx) error {
	doers := core.Doers{
		&validate{ctx: s.ctx},
		&buildURL{ctx: s.ctx},
		&response{ctx: s.ctx},
	}
	return doers.Do(doCtx)
}
