package controllers

import (
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type MemberController struct {
}

func (controller MemberController) Init(g *echo.Group) {
	g.GET("/my", controller.GetCurrentMember, middlewares.CheckPermission([]string{"*"}))
	g.GET("", controller.GetMembers, middlewares.CheckPermission([]string{"MANAGE_MEMBERS"}))
	g.GET("/:id", controller.GetMember, middlewares.CheckPermission([]string{"*"}))
	g.PUT("/:id/assign-roles", controller.AssignRole, middlewares.CheckPermission([]string{"MANAGE_MEMBERS"}))
}

func (MemberController) GetCurrentMember(ctx echo.Context) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	memberEntity, err := member.MemberService{}.GetMemberById(ctx.Request().Context(), userClaim.Id)

	memberInformation := dtos.CurrentMember{
		Id:          memberEntity.ID,
		Type:        memberEntity.Type,
		TypeName:    memberEntity.GetTypeName(),
		Name:        memberEntity.Name,
		Roles:       memberEntity.GetRoleNames(),
		Permissions: memberEntity.GetPermissionNames(),
	}
	return ctx.JSON(http.StatusOK, memberInformation)
}

func (MemberController) GetMembers(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)

	memberEntities, totalCount, err := member.MemberService{}.GetMembers(ctx.Request().Context(), nil, pageable)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var members = make([]dtos.MemberInformation, 0)
	for _, entity := range memberEntities {
		var roles = make([]dtos.MemberRole, 0)
		for _, memberRole := range entity.Roles {
			roles = append(roles, dtos.MemberRole{
				Id:   memberRole.ID,
				Name: memberRole.Name,
			})
		}
		members = append(members, dtos.MemberInformation{
			Id:          entity.ID,
			Type:        entity.Type,
			TypeName:    entity.GetTypeName(),
			Name:        entity.Name,
			MemberRoles: roles,
		})
	}

	pageResult := dtos.PageResult{
		Result:     members,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, pageResult)
}

func (MemberController) AssignRole(ctx echo.Context) error {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var assignRole dtos.MemberAssignRole
	if err := ctx.Bind(&assignRole); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := assignRole.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = member.MemberService{}.AssignRole(ctx.Request().Context(), uint(memberId), assignRole)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (MemberController) GetMember(ctx echo.Context) error {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	memberEntity, err := member.MemberService{}.GetMember(ctx.Request().Context(), uint(memberId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var roles = make([]dtos.MemberRole, 0)
	for _, memberRole := range memberEntity.Roles {
		roles = append(roles, dtos.MemberRole{
			Id:   memberRole.ID,
			Name: memberRole.Name,
		})
	}
	memberInformation := dtos.MemberInformation{
		Id:          memberEntity.ID,
		Type:        memberEntity.Type,
		TypeName:    memberEntity.GetTypeName(),
		Name:        memberEntity.Name,
		MemberRoles: roles,
	}

	return ctx.JSON(http.StatusOK, memberInformation)
}
