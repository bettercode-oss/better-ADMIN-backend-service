package app

import (
	"better-admin-backend-service/app/middlewares"
	"github.com/gin-contrib/cors"
	"time"
)

const AccessControlMaxAgeLimitHours = 24 // https://httptoolkit.com/blog/cache-your-cors

func (a *App) addGinMiddlewares() {
	a.gin.Use(cors.New(a.newCorsConfig()))
	a.gin.Use(middlewares.ErrorHandler)
	a.gin.Use(middlewares.JwtToken())
	a.gin.Use(middlewares.GORMDb(a.gormDB))
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