package helpers

import (
	"context"
	"gorm.io/gorm"
	"sync"
)

const ContextDBKey = "DB"

var (
	contextHelperOnce     sync.Once
	contextHelperInstance *contextHelper
)

func ContextHelper() *contextHelper {
	contextHelperOnce.Do(func() {
		contextHelperInstance = &contextHelper{}
	})

	return contextHelperInstance
}

type contextHelper struct {
}

func (contextHelper) GetDB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ContextDBKey)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*gorm.DB); ok {
		return db
	}
	panic("DB is not exist")
}

func (contextHelper) SetDB(ctx context.Context, gormDB *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextDBKey, gormDB)
}
