package routes

import "github.com/gin-gonic/gin"

type GinRoute interface {
	MapRoutes(routerGroup *gin.RouterGroup)
}
