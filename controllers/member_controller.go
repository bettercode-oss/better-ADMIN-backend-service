package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/factory"
	"better-admin-backend-service/domain/member"
	organiztion "better-admin-backend-service/domain/organization"
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
	g.POST("", controller.SignUpMember)
	g.GET("/my", controller.GetCurrentMember, middlewares.CheckPermission([]string{"*"}))
	g.GET("", controller.GetMembers, middlewares.CheckPermission([]string{domain.PermissionManageMembers}))
	g.GET("/:id", controller.GetMember, middlewares.CheckPermission([]string{"*"}))
	g.PUT("/:id/assign-roles", controller.AssignRole, middlewares.CheckPermission([]string{domain.PermissionManageMembers}))
	g.PUT("/:id/approved", controller.ApproveMember, middlewares.CheckPermission([]string{domain.PermissionManageMembers}))
}

func (MemberController) GetCurrentMember(ctx echo.Context) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	memberEntity, err := member.MemberService{}.GetMemberById(ctx.Request().Context(), userClaim.Id)
	memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx.Request().Context(), memberEntity)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	memberInformation := dtos.CurrentMember{
		Id:          memberEntity.ID,
		Type:        memberEntity.Type,
		TypeName:    memberEntity.GetTypeName(),
		Name:        memberEntity.Name,
		Roles:       memberAssignedAllRoleAndPermission.Roles,
		Permissions: memberAssignedAllRoleAndPermission.Permissions,
		Picture:     memberEntity.Picture,
	}

	return ctx.JSON(http.StatusOK, memberInformation)
}

func (MemberController) GetMembers(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)
	filters := map[string]interface{}{}

	if len(ctx.QueryParam("status")) > 0 {
		filters["status"] = ctx.QueryParam("status")
	}

	memberEntities, totalCount, err := member.MemberService{}.GetMembers(ctx.Request().Context(), filters, pageable)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	memberIds := make([]uint, 0)
	for _, entity := range memberEntities {
		memberIds = append(memberIds, entity.ID)
	}

	filters = map[string]interface{}{}
	filters["memberIds"] = memberIds
	organizationsOfMembers, err := organiztion.OrganizationService{}.GetAllOrganizations(ctx.Request().Context(), filters)

	var members = make([]dtos.MemberInformation, 0)
	for _, entity := range memberEntities {
		var roles = make([]dtos.MemberRole, 0)
		for _, memberRole := range entity.Roles {
			roles = append(roles, dtos.MemberRole{
				Id:   memberRole.ID,
				Name: memberRole.Name,
			})
		}
		memberInformation := dtos.MemberInformation{
			Id:          entity.ID,
			SignId:      entity.SignId,
			CandidateId: entity.GetCandidateId(),
			Type:        entity.Type,
			TypeName:    entity.GetTypeName(),
			Name:        entity.Name,
			MemberRoles: roles,
		}

		var memberOrganizations = make([]dtos.MemberOrganization, 0)
		for _, organizationsOfMember := range organizationsOfMembers {
			if organizationsOfMember.ExistMember(entity.ID) {
				memberOrganization := dtos.MemberOrganization{
					Id:   organizationsOfMember.ID,
					Name: organizationsOfMember.Name,
				}

				var memberOrganizationRoles = make([]dtos.MemberOrganizationRole, 0)
				for _, memberOrganizationRole := range organizationsOfMember.Roles {
					memberOrganizationRoles = append(memberOrganizationRoles, dtos.MemberOrganizationRole{
						Id:   memberOrganizationRole.ID,
						Name: memberOrganizationRole.Name,
					})
				}

				memberOrganization.Roles = memberOrganizationRoles
				memberOrganizations = append(memberOrganizations, memberOrganization)
			}
		}
		memberInformation.MemberOrganizations = memberOrganizations
		members = append(members, memberInformation)
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

func (MemberController) SignUpMember(ctx echo.Context) error {
	var memberSignUp dtos.MemberSignUp
	if err := ctx.Bind(&memberSignUp); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := memberSignUp.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := member.MemberService{}.SignUpMember(ctx.Request().Context(), memberSignUp)
	if err != nil {
		if err == domain.ErrDuplicated {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (MemberController) ApproveMember(ctx echo.Context) error {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = member.MemberService{}.ApproveMember(ctx.Request().Context(), uint(memberId))
	if err != nil {
		if err == domain.ErrAlreadyApproved {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
