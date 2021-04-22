package controllers

import (
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/middlewares"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	//"gopkg.in/testfixtures.v2"

	//"gopkg.in/testfixtures.v2"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	setUpEcho(setUpDatabase())
}

func setUpDatabase() *gorm.DB {
	fmt.Println("Set up database")
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	gormDB.AutoMigrate(&member.MemberEntity{}, &site.SettingEntity{})
	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	fmt.Println("Set up database test fixture")
	//fixtures, err := testfixtures.NewFolder(sqlDB, &testfixtures.SQLite{}, "../testdata/db_fixtures")
	//if err != nil {
	//	panic(err)
	//}
	//
	//testfixtures.SkipDatabaseNameCheck(true)
	//
	//if err := fixtures.Load(); err != nil {
	//	panic(err)
	//}

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),                      // You database connection
		testfixtures.Dialect("sqlite"),                    // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("../testdata/db_fixtures"), // the directory containing the YAML files
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)

	if err != nil {
		panic(err)
	}

	if err := fixtures.Load(); err != nil {
		panic(err)
	}

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
