package httpserver

import (
	"better-admin-backend-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGinEngine() *gin.Engine {
	e := gin.Default()

	e.SetTrustedProxies(nil) // https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	return e
}

func AddMiddlewares(e *gin.Engine, db *gorm.DB) *gin.Engine {
	e.Use(cors.New(newCorsConfig()))
	e.Use(middlewares.ErrorHandler)
	e.Use(middlewares.JwtToken())
	e.Use(middlewares.GORMDb(db))

	return e
}

func newCorsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	return corsConfig
}
