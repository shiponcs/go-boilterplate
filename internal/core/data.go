package core

import (
	"time"

	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/cache"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/models/entity"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/internal/stores"
)

// Ctx is the shared state base for every flow. Feature flows embed it in their
// own context struct (e.g. create.CreateCtx) and add flow-specific fields.
type Ctx struct {
	Srv    *services.SrvHolder
	Store  *stores.StoHolder
	Cache  cache.Cache
	Config *config.Config
	Log    *zap.Logger
	Now    time.Time

	Lang     string
	WidgetID uint

	Widget *entity.Widget
}
