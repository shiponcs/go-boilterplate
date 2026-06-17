package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/example/go-svc-boilerplate/internal/api"
	"github.com/example/go-svc-boilerplate/internal/api/handlers"
	"github.com/example/go-svc-boilerplate/internal/cache"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/conn"
	"github.com/example/go-svc-boilerplate/internal/core"
	"github.com/example/go-svc-boilerplate/internal/localization"
	"github.com/example/go-svc-boilerplate/internal/services"
	"github.com/example/go-svc-boilerplate/internal/stores"
	"github.com/example/go-svc-boilerplate/pkg"
)

// main is the single assembly point for the service. Every infra/domain layer
// exposes an fx module; top-level providers are listed directly here. To add a
// dependency, provide it in the relevant module (or here) and fx injects it by type.
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				pkg.NewLogger,
				config.LoadConfig,
			),
			// infra connections (db, redis, ...)
			conn.Module,
			// external service clients
			services.Module,
			// gorm-backed stores
			stores.Module,
			// redis-backed cache
			cache.Module,
		),
		fx.Provide(
			// handlers
			handlers.NewWidgetHandler,
			// core use-cases
			core.NewWidget,
			// http router
			api.SetupRoutes,
		),
		fx.Invoke(
			localization.Init,
			setupHTTPServer,
		),
	).Run()
}

// setupHTTPServer starts/stops the Gin server inside an fx lifecycle hook
// rather than a bare ListenAndServe, so shutdown is graceful.
func setupHTTPServer(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger, r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("starting HTTP server", zap.String("addr", srv.Addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("HTTP server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}
