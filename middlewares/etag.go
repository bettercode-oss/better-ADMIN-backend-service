package middlewares

import (
	"better-admin-backend-service/middlewares/internal"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpEtagCache(cacheControlMaxAgeSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := internal.NewResponseWriter(c)
		c.Writer = w
		defer w.Done(c)

		c.Next()

		eTag := getMD5Hash(w.Body().String())
		c.Writer.Header().Set("ETag", eTag)

		if len(c.Request.Header.Get("If-None-Match")) > 0 {
			if eTag == c.Request.Header.Get("If-None-Match") {
				c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%v", cacheControlMaxAgeSeconds))
				c.Status(http.StatusNotModified)
			}
		}
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
