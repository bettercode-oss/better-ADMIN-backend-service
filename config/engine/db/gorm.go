package db

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"os"
	"time"
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
)

func SetUpGormEngine() (*gorm.DB, error) {
	var dialector gorm.Dialector

	driver := os.Getenv(EnvDbDriver)

	if driver == "mysql" {
		if len(os.Getenv(EnvDbHost)) > 0 &&
			len(os.Getenv(EnvDbName)) > 0 &&
			len(os.Getenv(EnvDbUser)) > 0 &&
			len(os.Getenv(EnvDbPassword)) > 0 {
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				os.Getenv(EnvDbUser),
				os.Getenv(EnvDbPassword),
				os.Getenv(EnvDbHost),
				os.Getenv(EnvDbName))
			dialector = mysql.Open(dsn)
		} else {
			return nil, errors.New(fmt.Sprintf("%s, %s, %s  and %s environment variable are required.", EnvDbHost, EnvDbName, EnvDbUser, EnvDbPassword))
		}
	} else {
		// 기본적으로 DB는 sqlite
		dialector = sqlite.Open("account.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, errors.New("Database Connection Error")
	}

	if len(os.Getenv(EnvReplicaDbHost)) > 0 &&
		len(os.Getenv(EnvReplicaDbName)) > 0 &&
		len(os.Getenv(EnvReplicaDbUser)) > 0 &&
		len(os.Getenv(EnvReplicaDbPassword)) > 0 {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				os.Getenv(EnvReplicaDbUser),
				os.Getenv(EnvReplicaDbPassword),
				os.Getenv(EnvReplicaDbHost),
				os.Getenv(EnvReplicaDbName)))},
		}).SetConnMaxIdleTime(10).SetConnMaxLifetime(10 * time.Minute).SetMaxIdleConns(5).SetMaxOpenConns(10))
	}

	if err := initializeDatabase(db); err != nil {
		return nil, errors.Wrap(err, "Database initializeDatabase Error")
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}
