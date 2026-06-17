package core

import (
	"time"

	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/cache"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/dto"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/internal/stores"
)

// Widget is the domain use-case object handed to handlers. It owns the shared
// dependencies and is responsible for building the base flow Ctx; flows
// (create, get, ...) consume that Ctx. Mirror this for each new domain.
type Widget struct {
	cfg   *config.Config
	srv   *services.SrvHolder
	store *stores.StoHolder
	cache cache.Cache
	log   *zap.Logger
}

func NewWidget(
	cfg *config.Config,
	srv *services.SrvHolder,
	store *stores.StoHolder,
	cacheStore cache.Cache,
	log *zap.Logger,
) *Widget {
	return &Widget{
		cfg:   cfg,
		srv:   srv,
		store: store,
		cache: cacheStore,
		log:   log,
	}
}

// BaseCtx builds a flow context pre-populated with shared dependencies. Flows
// embed the returned Ctx and set their request-specific fields.
func (w *Widget) BaseCtx(lang string) Ctx {
	return Ctx{
		Srv:    w.srv,
		Store:  w.store,
		Cache:  w.cache,
		Config: w.cfg,
		Log:    w.log,
		Now:    time.Now(),
		Lang:   lang,
	}
}

// TransformWidget builds the API response from a stored widget. Response-shaping
// helpers like this belong on the use-case object, not in handlers.
func (w *Widget) TransformWidget(widget *entity.Widget) *dto.WidgetResponse {
	return &dto.WidgetResponse{
		ID:        widget.ID,
		Name:      widget.Name,
		Status:    widget.Status,
		Price:     widget.Price,
		CreatedAt: widget.CreatedAt.Unix(),
	}
}
