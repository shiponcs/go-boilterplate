package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger is a minimal structured request log middleware. Swap in your
// APM/tracing middleware here as the service grows.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		gin.DefaultWriter.Write([]byte(
			c.Request.Method + " " + c.Request.URL.Path + " " +
				time.Since(start).String() + "\n",
		))
	}
}
