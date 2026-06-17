package core

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/example/go-svc-boilerplate/internal/cache"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/stores"
	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// FetchWidget is a reusable step that loads a widget (cache first, then DB) into
// the flow context. Shared doers like this live in core and are embedded into
// any flow's Doers slice.
type FetchWidget struct {
	Ctx *Ctx
}

func (f *FetchWidget) Do(*DoCtx) error {
	widget, err := GetWidget(f.Ctx.Store.WidgetStore, f.Ctx.Cache, f.Ctx.WidgetID)
	if err != nil {
		return err
	}
	f.Ctx.Widget = widget
	return nil
}

// GetWidget reads through the cache to the store and back-fills the cache.
func GetWidget(widgetStore *stores.WidgetStore, cacheStore cache.Store, id uint) (*entity.Widget, error) {
	if cached, err := cacheStore.GetWidget(id); err == nil && cached != nil {
		return cached, nil
	}

	widget, err := widgetStore.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewNotFound("widget not found", "Widget not found")
		}
		return nil, err
	}

	_ = cacheStore.SetWidget(id, widget, 5*time.Minute)
	return widget, nil
}
