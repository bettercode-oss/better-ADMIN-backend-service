package rest

import (
	"better-admin-backend-service/app/middlewares"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type MemberController struct {
	routerGroup         *gin.RouterGroup
	rbacService         *services.RoleBasedAccessControlService
	memberService       *services.MemberService
	organizationService *services.OrganizationService
}

func NewMemberController(routerGroup *gin.RouterGroup,
	rbacService *services.RoleBasedAccessControlService,
	memberService *services.MemberService,
	organizationService *services.OrganizationService) *MemberController {

	return &MemberController{
		routerGroup:         routerGroup,
		rbacService:         rbacService,
		memberService:       memberService,
		organizationService: organizationService,
	}
}

func (c MemberController) MapRoutes() {
	route := c.routerGroup.Group("/members")

	route.POST("", c.signUpMember)
	route.GET("", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		c.getMembers)
	route.GET("/my", middlewares.PermissionChecker([]string{"*"}),
		c.getCurrentMember)
	route.GET("/:id", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		c.getMember)
	route.PUT("/:id/assign-roles", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		c.assignRole)
	route.PUT("/:id/approved", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		c.approveMember)
	route.PUT("/:id/rejected", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		c.rejectMember)
	route.GET("/search-filters", middlewares.PermissionChecker([]string{constants.PermissionManageMembers}),
		middlewares.HttpEtagCache(0),
		c.getSearchFilters)
}

func (c MemberController) signUpMember(ctx *gin.Context) {
	var memberSignUp dtos.MemberSignUp
	if err := ctx.BindJSON(&memberSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.memberService.SignUpMember(ctx.Request.Context(), memberSignUp)
	if err != nil {
		if err == errors.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c MemberController) getCurrentMember(ctx *gin.Context) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	memberEntity, err := c.memberService.GetMemberById(ctx.Request.Context(), userClaim.Id)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	memberAssignedAllRoleAndPermission, err := c.organizationService.GetMemberAssignedAllRoleAndPermission(ctx.Request.Context(), memberEntity)
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

func (c MemberController) getMembers(ctx *gin.Context) {
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

	memberEntities, totalCount, err := c.memberService.GetMembers(ctx.Request.Context(), filters, pageable)
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
	organizationsOfMembers, err := c.organizationService.GetAllOrganizations(ctx.Request.Context(), filters)
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

func (c MemberController) getMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	memberEntity, err := c.memberService.GetMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == errors.ErrNotFound {
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

func (c MemberController) assignRole(ctx *gin.Context) {
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

	err = c.memberService.AssignRole(ctx.Request.Context(), uint(memberId), assignRole)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c MemberController) approveMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.memberService.ApproveMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == errors.ErrAlreadyApproved {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c MemberController) getSearchFilters(ctx *gin.Context) {
	filters := make([]dtos.SearchFilter, 0)

	memberTypeSearchFilter := dtos.SearchFilter{
		Name: "type",
		Filters: []dtos.Filter{
			{
				Text:  constants.TypeMemberSiteName,
				Value: constants.TypeMemberSite,
			},
			{
				Text:  constants.TypeMemberDoorayName,
				Value: constants.TypeMemberDooray,
			},
			{
				Text:  constants.TypeMemberGoogleName,
				Value: constants.TypeMemberGoogle,
			},
		},
	}
	filters = append(filters, memberTypeSearchFilter)

	allRoles, _, err := c.rbacService.GetRoles(ctx.Request.Context(), nil, dtos.Pageable{Page: 0})
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

func (c MemberController) rejectMember(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.memberService.RejectMember(ctx.Request.Context(), uint(memberId))
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
