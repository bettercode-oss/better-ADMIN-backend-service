package controllers

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/config/engine/httpserver"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	gormDB *gorm.DB
	ginApp *gin.Engine
)

func init() {
	setUpTestGinServer()
}

func setUpTestGinServer() {
	config.InitConfig("../config/config.json")
	ginApp = httpserver.NewGinEngine()
	httpserver.AddMiddlewares(ginApp, setUpTestDatabase())
	AddRoutes(ginApp)
}

func setUpTestDatabase() *gorm.DB {
	fmt.Println("Set up database")
	db, err := gorm.Open(sqlite.Open("file::memory:?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}
	gormDB = db
	return gormDB
}

func generateTestJWT(claim map[string]interface{}, duration time.Duration) (string, error) {
	token := jwt.MapClaims{}
	for key, value := range claim {
		token[key] = value
	}

	token["exp"] = time.Now().Add(duration).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(config.Config.JwtSecret))
}
