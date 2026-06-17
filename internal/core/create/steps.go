package create

import (
	"strings"

	"github.com/example/go-svc-boilerplate/internal/cnst"
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/localization"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// validate rejects bad input before any work is done.
type validate struct {
	ctx *CreateCtx
}

func (v *validate) Do(*core.DoCtx) error {
	if strings.TrimSpace(v.ctx.ReqDt.Name) == "" {
		msg := localization.GetMessage("widget_invalid_name", v.ctx.Lang)
		return errs.NewBadReq("empty widget name", msg)
	}
	return nil
}

// quote computes the price via the pricing service and stores it on the ctx.
type quote struct {
	ctx *CreateCtx
}

func (q *quote) Do(*core.DoCtx) error {
	pricing, err := q.ctx.Srv.PricingService.Quote(q.ctx.ReqDt.Units)
	if err != nil {
		return errs.NewInternalServer("failed to quote price: "+err.Error(), "")
	}
	q.ctx.pricing = pricing
	return nil
}

// persist writes the widget to the store and back-fills the flow context.
type persist struct {
	ctx *CreateCtx
}

func (p *persist) Do(*core.DoCtx) error {
	widget := &entity.Widget{
		Name:   p.ctx.ReqDt.Name,
		Status: cnst.WidgetStatusActive,
		Price:  p.ctx.pricing.Total,
	}
	if err := p.ctx.Store.WidgetStore.Create(widget); err != nil {
		return errs.NewInternalServer("failed to persist widget: "+err.Error(), "")
	}
	p.ctx.Widget = widget
	return nil
}

// response builds the API payload from the persisted widget.
type response struct {
	ctx *CreateCtx
}

func (r *response) Do(*core.DoCtx) error {
	w := r.ctx.Widget
	r.ctx.Resp = &dto.WidgetResponse{
		ID:        w.ID,
		Name:      w.Name,
		Status:    w.Status,
		Price:     w.Price,
		CreatedAt: w.CreatedAt.Unix(),
	}
	return nil
}
