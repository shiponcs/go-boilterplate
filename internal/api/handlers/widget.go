package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/core/create"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// WidgetHandler is the HTTP entry point for the widget domain. It binds the
// request, builds the flow context from the use-case object, runs the flow, and
// serializes the result. No business logic lives here.
type WidgetHandler struct {
	widget *core.Widget
}

func NewWidgetHandler(widget *core.Widget) *WidgetHandler {
	return &WidgetHandler{widget: widget}
}

func (h *WidgetHandler) CreateWidget(c *gin.Context) {
	var req dto.CreateWidgetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ServeErr(c, errs.NewBadReq("invalid request body: "+err.Error(), "Invalid request"))
		return
	}

	ctx := &create.CreateCtx{
		Ctx:   h.widget.BaseCtx(lang(c)),
		ReqDt: &req,
	}
	if err := create.New(ctx).Do(&core.DoCtx{}); err != nil {
		ServeErr(c, err)
		return
	}

	ServeData(c, ctx.Resp, 201)
}

func (h *WidgetHandler) GetWidget(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ServeErr(c, errs.NewBadReq("invalid widget id", "Invalid widget id"))
		return
	}

	base := h.widget.BaseCtx(lang(c))
	base.WidgetID = uint(id)

	fetch := &core.FetchWidget{Ctx: &base}
	if err := fetch.Do(&core.DoCtx{}); err != nil {
		ServeErr(c, err)
		return
	}

	ServeData(c, h.widget.TransformWidget(base.Widget))
}

// lang resolves the request language, defaulting to English.
func lang(c *gin.Context) string {
	if l := c.GetHeader("Accept-Language"); l != "" {
		return l
	}
	return "en"
}
