package app

import (
	"better-admin-backend-service/app/db"
	"better-admin-backend-service/app/routes"
	"better-admin-backend-service/http/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
)

type App struct {
	gormDB            *gorm.DB
	webSocketUpgrader websocket.Upgrader
	gin               *gin.Engine
	router            routes.GinRoute
	dbConnector       db.DatabaseConnector
}

func NewApp(router routes.GinRoute, dbConnector db.DatabaseConnector) *App {
	g := gin.Default()
	g.SetTrustedProxies(nil) // https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	return &App{gin: g, webSocketUpgrader: upgrader, router: router, dbConnector: dbConnector}
}

func (a *App) SetUp() error {
	db, err := a.dbConnector.Connect()
	if err != nil {
		return err
	}
	a.gormDB = db

	if err := a.migrateDatabase(); err != nil {
		return err
	}

	a.gin.GET("/ws/:id", ws.WebSocketHandler(a.webSocketUpgrader))

	// Liveness Probe
	a.gin.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	a.addGinMiddlewares()
	a.router.MapRoutes(a.gin.Group("/api"))
	return nil
}

func (a *App) Run() error {
	a.SetUp()
	sqlDB, err := a.gormDB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	a.gin.Run(":2016")
	return nil
}

func (a App) GetGin() *gin.Engine {
	return a.gin
}

func (a App) GetDB() *gorm.DB {
	return a.gormDB
}
