package testdb

import (
	"better-admin-backend-service/domain/member/entity"
	entity2 "better-admin-backend-service/domain/organization/entity"
	entity5 "better-admin-backend-service/domain/rbac/entity"
	entity4 "better-admin-backend-service/domain/site/entity"
	entity3 "better-admin-backend-service/domain/webhook/entity"
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/gorm"
)

type DatabaseFixture struct {
}

func (DatabaseFixture) SetUpDefault(gormDB *gorm.DB) {
	fmt.Println("Set up database test fixture")
	// TODO 중복 없애기 필요
	gormDB.AutoMigrate(&entity.MemberEntity{}, &entity4.SettingEntity{}, &entity5.PermissionEntity{}, &entity5.RoleEntity{},
		&entity2.OrganizationEntity{}, &entity3.WebHookEntity{}, &entity3.WebHookMessageEntity{})

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),                               // You database connection
		testfixtures.Dialect("sqlite"),                             // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("../testdata/testdb/data_fixtures"), // the directory containing the YAML files
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
