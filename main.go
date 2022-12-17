package main

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/config"
	"better-admin-backend-service/config/engine/db"
	"better-admin-backend-service/config/engine/httpserver"
	"better-admin-backend-service/controllers"
	"better-admin-backend-service/helpers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	filename "github.com/keepeye/logrus-filename"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
)

const (
	EnvJwtSecret = "JWT_SECRET"
)

var (
	gormDB   *gorm.DB
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
)

func init() {
	config.InitConfig("config/config.json")
	setUpJwtSecret()
	setUpLogFormatter()
}

func setUpJwtSecret() {
	if len(os.Getenv(EnvJwtSecret)) > 0 {
		config.Config.JwtSecret = os.Getenv(EnvJwtSecret)
	}
}

func setUpLogFormatter() {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	log.AddHook(filenameHook)
	log.SetFormatter(&log.JSONFormatter{DisableHTMLEscape: true})
}

func main() {
	gormDB, err := db.SetUpGormEngine()
	if err != nil {
		panic(err.Error())
	}

	g := httpserver.NewGinEngine()
	g.GET("/ws/:id", connectWebSocket)

	httpserver.AddMiddlewares(g, gormDB)
	controllers.AddRoutes(g)
	g.Run(":2016")
}

func connectWebSocket(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	webSocketId := ctx.Param("id")
	adapters.WebSocketAdapter().AddConnection(webSocketId, ws)

	ctx.Status(http.StatusOK)
}
