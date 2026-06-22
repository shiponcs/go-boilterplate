package api

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/api/handlers"
	"github.com/example/go-svc-boilerplate/internal/api/middleware"
	"github.com/example/go-svc-boilerplate/internal/config"
	"github.com/example/go-svc-boilerplate/internal/core"
)

// SetupRoutes builds the Gin engine and registers routes. It is the single
// place that maps URLs to handler methods.
func SetupRoutes(wh *handlers.WidgetHandler, ah *handlers.AuthHandler, auth *core.Auth, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger())

	pprof.Register(r)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// requireAuth guards routes that need an authenticated session. Apply it to
	// any protected route (e.g. the widget routes) as the service grows.
	requireAuth := middleware.RequireAuth(auth, cfg)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/widgets", wh.CreateWidget)
		v1.GET("/widgets/:id", wh.GetWidget)

		v1.GET("/auth/signup", ah.Signup)
		v1.GET("/auth/callback", ah.Callback)
		v1.GET("/auth/logout", ah.Logout)
		v1.GET("/auth/me", requireAuth, ah.Me)
	}

	return r
}
