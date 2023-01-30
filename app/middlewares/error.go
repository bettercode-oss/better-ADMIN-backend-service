package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if c.Writer.Status() >= 500 {
		for _, err := range c.Errors {
			log.Errorf("%+v", err.Err)
		}
		c.Status(http.StatusInternalServerError)
	}
}
