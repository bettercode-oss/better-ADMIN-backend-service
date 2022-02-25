package controllers

import (
	"better-admin-backend-service/domain/logging"
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/menu"
	"better-admin-backend-service/domain/organization"
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/domain/webhook"
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseFixture struct {
}

func (DatabaseFixture) setUpDefault() {
	fmt.Println("Set up database test fixture")
	gormDB.AutoMigrate(&member.MemberEntity{}, &site.SettingEntity{}, &rbac.PermissionEntity{}, &rbac.RoleEntity{},
		&organization.OrganizationEntity{}, &webhook.WebHookEntity{}, &webhook.WebHookMessageEntity{}, &menu.MenuEntity{},
		&logging.MemberAccessLogEntity{})

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

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
	fmt.Println("End of database test fixture")
}
