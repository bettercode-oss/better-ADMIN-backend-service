package controllers

import (
	"better-admin-backend-service/middlewares"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	//"gopkg.in/testfixtures.v2"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	gormDB           *gorm.DB
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	setUpEcho(setUpDatabase())
}

func setUpDatabase() *gorm.DB {
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

func setUpEcho(gormDB *gorm.DB) {
	fmt.Println("Set up echo")
	echoApp = echo.New()
	echoApp.Validator = &CustomValidator{validator: validator.New()}

	db := middlewares.GORMDb(gormDB)

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return db(handlerFunc)(c)
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
