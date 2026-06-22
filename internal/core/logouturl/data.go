package logouturl

import (
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
)

// LogoutURLCtx embeds the shared core.Ctx and carries input/output for
// building the AuthKit logout URL.
type LogoutURLCtx struct {
	core.Ctx

	SessionID string
	ReturnTo  string

	url  string
	Resp *dto.LogoutURLResponse
}
