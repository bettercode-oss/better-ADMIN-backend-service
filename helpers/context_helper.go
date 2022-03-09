package helpers

import (
	"better-admin-backend-service/security"
	"context"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
	"sync"
)

const ContextDBKey = "DB"
const ContextUserClaimKey = "userClaim"

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

func (contextHelper) SetUserClaim(ctx context.Context, userClaim *security.UserClaim) context.Context {
	return context.WithValue(ctx, ContextUserClaimKey, userClaim)
}

func (contextHelper) GetUserClaim(ctx context.Context) (*security.UserClaim, error) {
	v := ctx.Value(ContextUserClaimKey)
	if v == nil {
		return nil, errors.New("UserClaim is not exist")
	}
	if userClaim, ok := v.(*security.UserClaim); ok {
		return userClaim, nil
	}
	return nil, errors.New("UserClaim is not exist")
}
