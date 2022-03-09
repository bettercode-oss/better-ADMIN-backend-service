package middlewares

import (
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil {
				log.Error(err.Error(), err.(*errors.Error).ErrorStack())
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return err
		}
	}
}
