package create

import (
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/value"
)

// CreateCtx embeds the shared core.Ctx and adds fields specific to the create
// flow. Steps read/write these fields as they run.
type CreateCtx struct {
	core.Ctx

	ReqDt *dto.CreateWidgetReq
	Resp  *dto.WidgetResponse

	pricing *value.CalculatedPrice
}
