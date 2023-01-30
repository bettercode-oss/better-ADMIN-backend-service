package testdb

import (
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/gorm"
)

type DatabaseFixture struct {
}

func (DatabaseFixture) SetUpDefault(gormDB *gorm.DB) {
	fmt.Println("Set up database test fixture")
	sqlDB, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),                                  // You database connection
		testfixtures.Dialect("sqlite"),                                // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory("../../testdata/testdb/data_fixtures"), // the directory containing the YAML files
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
