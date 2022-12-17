package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/member/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type MemberController struct {
}

func (controller MemberController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/members")

	route.POST("", controller.signUpMember)
	route.GET("", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		controller.getMembers)
	route.GET("/my", middlewares.PermissionChecker([]string{"*"}),
		controller.getCurrentMember)
	route.GET("/:id", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		controller.getMember)
	route.PUT("/:id/assign-roles", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		controller.assignRole)
	route.PUT("/:id/approved", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		controller.approveMember)
	route.PUT("/:id/rejected", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		controller.rejectMember)
	route.GET("/search-filters", middlewares.PermissionChecker([]string{domain.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		controller.getSearchFilters)
}

func (MemberController) signUpMember(ctx *gin.Context) {
	var memberSignUp dtos.MemberSignUp
	if err := ctx.BindJSON(&memberSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := services.MemberService{}.SignUpMember(ctx.Request.Context(), memberSignUp)
	if err != nil {
		if err == domain.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (MemberController) getCurrentMember(ctx *gin.Context) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	memberEntity, err := services.MemberService{}.GetMemberById(ctx.Request.Context(), userClaim.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	memberAssignedAllRoleAndPermission, err := services.OrganizationService{}.GetMemberAssignedAllRoleAndPermission(ctx.Request.Context(), memberEntity)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, memberInformation)
}

func (MemberController) getMembers(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)
	filters := map[string]interface{}{}

	if len(ctx.Query("status")) > 0 {
		filters["status"] = ctx.Query("status")
	}

	if len(ctx.Query("name")) > 0 {
		filters["name"] = ctx.Query("name")
	}

	if len(ctx.Query("types")) > 0 {
		filters["types"] = strings.Split(ctx.Query("types"), ",")
	}

	if len(ctx.Query("roleIds")) > 0 {
		filters["roleIds"] = strings.Split(ctx.Query("roleIds"), ",")
	}

	memberEntities, totalCount, err := services.MemberService{}.GetMembers(ctx.Request.Context(), filters, pageable)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	memberIds := make([]uint, 0)
	for _, entity := range memberEntities {
		memberIds = append(memberIds, entity.ID)
	}

	filters = map[string]interface{}{}
	filters["memberIds"] = memberIds
	organizationsOfMembers, err := services.OrganizationService{}.GetAllOrganizations(ctx.Request.Context(), filters)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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
		memberInformation := dtos.MemberInformation{
			Id:           entity.ID,
			SignId:       entity.SignId,
			CandidateId:  entity.GetCandidateId(),
			Type:         entity.Type,
			TypeName:     entity.GetTypeName(),
			Name:         entity.Name,
			MemberRoles:  roles,
			CreatedAt:    entity.CreatedAt,
			LastAccessAt: entity.LastAccessAt,
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

	ctx.JSON(http.StatusOK, pageResult)
}

func (MemberController) getMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	memberEntity, err := services.MemberService{}.GetMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, memberInformation)
}

func (MemberController) assignRole(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var assignRole dtos.MemberAssignRole
	if err := ctx.BindJSON(&assignRole); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.MemberService{}.AssignRole(ctx.Request.Context(), uint(memberId), assignRole)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (MemberController) approveMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.MemberService{}.ApproveMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == domain.ErrAlreadyApproved {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (MemberController) getSearchFilters(ctx *gin.Context) {
	filters := make([]dtos.SearchFilter, 0)

	memberTypeSearchFilter := dtos.SearchFilter{
		Name: "type",
		Filters: []dtos.Filter{
			{
				Text:  entity.TypeMemberSiteName,
				Value: entity.TypeMemberSite,
			},
			{
				Text:  entity.TypeMemberDoorayName,
				Value: entity.TypeMemberDooray,
			},
			{
				Text:  entity.TypeMemberGoogleName,
				Value: entity.TypeMemberGoogle,
			},
		},
	}
	filters = append(filters, memberTypeSearchFilter)

	allRoles, _, err := services.RoleBasedAccessControlService{}.GetRoles(ctx.Request.Context(), nil, dtos.Pageable{Page: 0})
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	roleSearchFilter := dtos.SearchFilter{
		Name: "role",
	}
	roleFilters := make([]dtos.Filter, 0)
	for _, role := range allRoles {
		roleFilters = append(roleFilters, dtos.Filter{
			Text:  role.Name,
			Value: strconv.FormatUint(uint64(role.ID), 10),
		})
	}
	roleSearchFilter.Filters = roleFilters
	filters = append(filters, roleSearchFilter)

	ctx.JSON(http.StatusOK, filters)
}

func (MemberController) rejectMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.MemberService{}.RejectMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
