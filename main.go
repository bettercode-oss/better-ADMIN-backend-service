package main

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/config"
	"better-admin-backend-service/controllers"
	"better-admin-backend-service/middlewares"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	filename "github.com/keepeye/logrus-filename"
	"github.com/labstack/echo"
	echomiddleware "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"net/http"
	"os"
	"time"
)

const (
	Version = "0.0.1"
	website = "https://www.bettercode.kr"
	banner  = `
  __        __   __                    ___    ___    __  ___   ____   _  __
  / /  ___  / /_ / /_ ___   ____       / _ |  / _ \  /  |/  /  /  _/  / |/ /
 / _ \/ -_)/ __// __// -_) / __/      / __ | / // / / /|_/ /  _/ /   /    / 
/_.__/\__/ \__/ \__/ \__/ /_/        /_/ |_|/____/ /_/  /_/  /___/  /_/|_/  
`
)

const (
	EnvDbDriver          = "DB_DRIVER"
	EnvDbHost            = "DB_HOST"
	EnvDbName            = "DB_NAME"
	EnvDbUser            = "DB_USER"
	EnvDbPassword        = "DB_PASSWORD"
	EnvReplicaDbHost     = "REPLICA_DB_HOST"
	EnvReplicaDbName     = "REPLICA_DB_NAME"
	EnvReplicaDbUser     = "REPLICA_DB_USER"
	EnvReplicaDbPassword = "REPLICA_DB_PASSWORD"
	EnvJwtSecret         = "JWT_SECRET"
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
	setUpDatabase()
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

func setUpDatabase() {
	var dialector gorm.Dialector

	driver := os.Getenv(EnvDbDriver)

	if driver == "mysql" {
		if len(os.Getenv(EnvDbHost)) > 0 &&
			len(os.Getenv(EnvDbName)) > 0 &&
			len(os.Getenv(EnvDbUser)) > 0 &&
			len(os.Getenv(EnvDbPassword)) > 0 {
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
				os.Getenv(EnvDbUser),
				os.Getenv(EnvDbPassword),
				os.Getenv(EnvDbHost),
				os.Getenv(EnvDbName))
			dialector = mysql.Open(dsn)
		} else {
			panic(fmt.Sprintf("%s, %s, %s  and %s environment variable are required.", EnvDbHost, EnvDbName, EnvDbUser, EnvDbPassword))
		}
	} else {
		// 기본적으로 DB는 sqlite
		dialector = sqlite.Open("account.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Database Connection Error")
	}

	if len(os.Getenv(EnvReplicaDbHost)) > 0 &&
		len(os.Getenv(EnvReplicaDbName)) > 0 &&
		len(os.Getenv(EnvReplicaDbUser)) > 0 &&
		len(os.Getenv(EnvReplicaDbPassword)) > 0 {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
				os.Getenv(EnvReplicaDbUser),
				os.Getenv(EnvReplicaDbPassword),
				os.Getenv(EnvReplicaDbHost),
				os.Getenv(EnvReplicaDbName)))},
		}).SetConnMaxIdleTime(10).SetConnMaxLifetime(10 * time.Minute).SetMaxIdleConns(5).SetMaxOpenConns(10))
	}

	if err := initializeDatabase(db); err != nil {
		panic(fmt.Sprintf("Database initializeDatabase Error : %s", err.Error()))
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	gormDB = db
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Pre(echomiddleware.RemoveTrailingSlash())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowCredentials: true,
	}))

	e.Use(middlewares.JwtToken())
	e.Use(middlewares.GORMDb(gormDB))
	e.HideBanner = true

	e.GET("/ws/:id", connectWebSocket)
	controllers.AuthController{}.Init(e.Group("/api/auth"))
	controllers.SiteController{}.Init(e.Group("/api/site"))
	controllers.MemberController{}.Init(e.Group("/api/members"))
	controllers.AccessControlController{}.Init(e.Group("/api/access-control"))
	controllers.OrganizationController{}.Init(e.Group("/api/organizations"))
	controllers.WebHookController{}.Init(e.Group("/api/web-hooks"))

	color.Println(banner, color.Red("v"+Version), color.Blue(website))
	e.Start(":2016")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func connectWebSocket(ctx echo.Context) error {
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	webSocketId := ctx.Param("id")
	adapters.WebSocketAdapter().AddConnection(webSocketId, ws)

	return nil
}
