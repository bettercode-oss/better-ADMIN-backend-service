package middlewares

import (
	"better-admin-backend-service/helpers"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GORMDb(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			switch req.Method {
			case "POST", "PUT", "DELETE", "PATCH":
				tx := db.Begin()
				defer func() {
					if r := recover(); r != nil {
						tx.Rollback()
					}
				}()

				if err := tx.Error; err != nil {
					return err
				}

				ctx = helpers.ContextHelper().SetDB(ctx, tx)
				c.SetRequest(req.WithContext(ctx))

				if err := next(c); err != nil {
					tx.Rollback()
					return err
				}
				if c.Response().Status >= 500 {
					if err := tx.Rollback().Error; err != nil {
						log.Error("database rollback error", err.Error())
					}
					return nil
				}
				if err := tx.Commit().Error; err != nil {
					log.Error("database commit error", err.Error())
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				ctx = helpers.ContextHelper().SetDB(ctx, db)
				c.SetRequest(req.WithContext(ctx))
				return next(c)
			}

			return nil
		}
	}
}
