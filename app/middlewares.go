package app

import (
	"better-admin-backend-service/app/middlewares"
	xss "github.com/bettercode-oss/gin-middleware-xss"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"
)

const AccessControlMaxAgeLimitHours = 24 // https://httptoolkit.com/blog/cache-your-cors

func (a *App) addGinMiddlewares() {
	a.gin.Use(cors.New(a.newCorsConfig()))
	a.gin.Use(middlewares.ErrorHandler)
	a.gin.Use(middlewares.JwtToken())
	a.gin.Use(middlewares.GORMDb(a.gormDB))
	a.gin.Use(xss.Sanitizer(xss.Config{
		UrlsToExclude:     []string{"/api/auth", "/api/auth/dooray"},
		TargetHttpMethods: []string{http.MethodPost, http.MethodPut}}))
}

func (a *App) newCorsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.MaxAge = AccessControlMaxAgeLimitHours * time.Hour
	return corsConfig
}
