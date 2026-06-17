package api

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/internal/api/handlers"
	"github.com/example/go-svc-boilerplate/internal/api/middleware"
)

// SetupRoutes builds the Gin engine and registers routes. It is the single
// place that maps URLs to handler methods.
func SetupRoutes(wh *handlers.WidgetHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger())

	pprof.Register(r)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/widgets", wh.CreateWidget)
		v1.GET("/widgets/:id", wh.GetWidget)
	}

	return r
}
