package helpers

import (
	"better-admin-backend-service/dtos"
	"gorm.io/gorm"
	"sync"
)

var (
	gormHelperOnce     sync.Once
	gormHelperInstance *gormHelper
)

func GormHelper() *gormHelper {
	gormHelperOnce.Do(func() {
		gormHelperInstance = &gormHelper{}
	})

	return gormHelperInstance
}

type gormHelper struct {
}

func (gormHelper) Pageable(pageable dtos.Pageable) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageable.Page > 0 {
			return db.Limit(pageable.PageSize).Offset(pageable.GetOffset())
		}
		return db
	}
}
