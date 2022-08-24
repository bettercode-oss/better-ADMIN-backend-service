package controllers

import (
	entity6 "better-admin-backend-service/domain/logging/entity"
	"better-admin-backend-service/domain/member/entity"
	entity2 "better-admin-backend-service/domain/organization/entity"
	entity5 "better-admin-backend-service/domain/rbac/entity"
	entity4 "better-admin-backend-service/domain/site/entity"
	entity3 "better-admin-backend-service/domain/webhook/entity"
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseFixture struct {
}

func (DatabaseFixture) setUpDefault() {
	fmt.Println("Set up database test fixture")
	gormDB.AutoMigrate(&entity.MemberEntity{}, &entity4.SettingEntity{}, &entity5.PermissionEntity{}, &entity5.RoleEntity{},
		&entity2.OrganizationEntity{}, &entity3.WebHookEntity{}, &entity3.WebHookMessageEntity{},
		&entity6.MemberAccessLogEntity{})

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
