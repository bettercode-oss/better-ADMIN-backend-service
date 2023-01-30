package middlewares

import (
	"better-admin-backend-service/helpers"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GORMDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		ctx := req.Context()

		switch req.Method {
		case "POST", "PUT", "DELETE", "PATCH":
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					panic(r) // 상위 middleware(Recover) 에서 Panic 을 처리하도록 함.
				}
			}()

			if err := tx.Error; err != nil {
				c.Error(errors.Wrap(err, "DB Tx Begin error"))
				c.Abort()
				return
			}
			c.Request = c.Request.WithContext(helpers.ContextHelper().SetDB(ctx, tx))

			c.Next()

			if len(c.Errors) > 0 {
				if err := tx.Rollback().Error; err != nil {
					c.Error(errors.Wrap(err, "database rollback error"))
				}
				c.Abort()
				return
			}

			if err := tx.Commit().Error; err != nil {
				c.Error(errors.Wrap(err, "database rollback error"))
				c.Abort()
				return
			}
		default:
			c.Request = c.Request.WithContext(helpers.ContextHelper().SetDB(ctx, db))
			c.Next()
		}
	}
}
