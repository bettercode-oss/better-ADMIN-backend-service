package testserver

import (
	"better-admin-backend-service/app"
	"better-admin-backend-service/app/routes"
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewTestAppServer(router routes.GinRoute) *app.App {
	opaRego, err := rego.New(rego.Query("data.rest.allowed"), rego.Load([]string{
		"../../authorization/rest/policy.rego", "../../authorization/rest/data.json",
	}, nil)).PrepareForEval(context.TODO())

	if err != nil {
		panic(err)
	}

	testApp := app.NewApp(router, TestDbConnector{}, &opaRego)
	err = testApp.SetUp()
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
