package testserver

import (
	"better-admin-backend-service/app"
	"better-admin-backend-service/app/routes"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestAppServer(router routes.GinRoute) *app.App {
	testApp := app.NewApp(router, TestDbConnector{})
	err := testApp.SetUp()
	if err != nil {
		panic(err)
	}

	return testApp
}

type TestDbConnector struct {
}

func (TestDbConnector) Connect() (*gorm.DB, error) {
	fmt.Println("Set up database")
	return gorm.Open(sqlite.Open("file::memory:?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
