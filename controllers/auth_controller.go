package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/auth"
	"better-admin-backend-service/dtos"
	"github.com/labstack/echo"
	"net/http"
)

type AuthController struct {
}

func (controller AuthController) Init(g *echo.Group) {
	g.POST("", controller.AuthWithSignIdPassword)
}

func (AuthController) AuthWithSignIdPassword(ctx echo.Context) error {
	var memberSignIn dtos.MemberSignIn

	if err := ctx.Bind(&memberSignIn); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := memberSignIn.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	token, err := auth.AuthService{}.AuthWithSignIdPassword(ctx.Request().Context(), memberSignIn)
	if err != nil {
		if err == domain.ErrNotFound || err == domain.ErrAuthentication {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, token)
}
