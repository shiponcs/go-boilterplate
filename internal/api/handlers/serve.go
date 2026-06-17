package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/go-svc-boilerplate/pkg/errs"
)

// ServeData writes a JSON body with an optional status code (defaults to 200).
func ServeData(c *gin.Context, data any, codes ...int) {
	status := http.StatusOK
	if len(codes) > 0 && codes[0] >= 100 && codes[0] < 600 {
		status = codes[0]
	}
	c.JSON(status, data)
}

// ServeErr inspects an error and writes the appropriate response. Custom errs
// errors render their own status/body; everything else becomes a generic 500.
func ServeErr(c *gin.Context, e error) {
	log.Printf("serving http error: %v", e)

	var customErr interface {
		errs.StatusCoder
		errs.HTTPFormatter
	}
	if errors.As(e, &customErr) {
		c.JSON(customErr.StatusCode(), customErr.HTTPFormat())
		return
	}

	internal(c)
}

func internal(c *gin.Context) {
	err := errs.NewInternalServer("unhandled internal error", "something went wrong!")
	formatter := err.(errs.HTTPFormatter)
	statusCode := err.(errs.StatusCoder)
	c.JSON(statusCode.StatusCode(), formatter.HTTPFormat())
}
