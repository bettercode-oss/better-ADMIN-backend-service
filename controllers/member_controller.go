package controllers

import (
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
)

type MemberController struct {
}

func (controller MemberController) Init(g *echo.Group) {
	g.GET("/my", controller.GetCurrentMember, middlewares.CheckAuth())
}

func (MemberController) GetCurrentMember(ctx echo.Context) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	memberEntity, err := member.MemberService{}.GetMemberById(ctx.Request().Context(), userClaim.Id)

	memberInformation := dtos.MemberInformation{
		Id:   memberEntity.ID,
		Name: memberEntity.Name,
	}
	return ctx.JSON(http.StatusOK, memberInformation)
}
