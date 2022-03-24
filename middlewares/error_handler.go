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
				if _, ok := err.(*errors.Error); ok {
					log.Error(err.Error(), err.(*errors.Error).ErrorStack())
				} else {
					log.Error(err.Error())
				}
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return err
		}
	}
}
