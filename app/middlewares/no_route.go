package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var ginRouteMap map[string]map[string]bool
var once sync.Once

func getRouteMapInstance(routes gin.RoutesInfo) map[string]map[string]bool {
	once.Do(func() {
		routeMap := map[string]map[string]bool{}
		for _, r := range routes {
			if routeMap[r.Path] != nil {
				routeMap[r.Path][r.Method] = true
			} else {
				routeMap[r.Path] = map[string]bool{
					r.Method: true,
				}
			}
		}
		ginRouteMap = routeMap
	})
	return ginRouteMap
}

func NoRoute(ginEngine *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		routeMap := getRouteMapInstance(ginEngine.Routes())
		if routeMap[c.FullPath()] == nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if routeMap[c.FullPath()][c.Request.Method] == false {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		c.Next()
	}
}
