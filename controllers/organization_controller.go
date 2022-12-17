package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/factory"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type OrganizationController struct {
}

func (controller OrganizationController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/organizations")

	route.POST("", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.createOrganization)
	route.GET("", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		middlewares.HttpEtagCache(0),
		controller.getOrganizations)
	route.GET("/:organizationId", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		middlewares.HttpEtagCache(0),
		controller.getOrganization)
	route.PUT("/:organizationId/name", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.changeOrganizationName)
	route.PUT("/:organizationId/change-position", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.changePosition)
	route.PUT("/:organizationId/assign-roles", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.assignRoles)
	route.PUT("/:organizationId/assign-members", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.assignMembers)
	route.DELETE("/:organizationId", middlewares.PermissionChecker([]string{domain.PermissionManageOrganization}),
		controller.deleteOrganization)
}

func (controller OrganizationController) createOrganization(ctx *gin.Context) {
	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := services.OrganizationService{}.CreateOrganization(ctx.Request.Context(), organizationInformation)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (controller OrganizationController) getOrganizations(ctx *gin.Context) {
	allOfOrganizations, err := services.OrganizationService{}.GetAllOrganizations(ctx.Request.Context(), nil)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	organizations := make([]dtos.OrganizationInformation, 0)
	for _, entity := range allOfOrganizations {
		if entity.ParentOrganizationID == nil {
			organizationInformation := factory.NewOrganizationInformationFromEntity(entity)
			organizations = append(organizations, organizationInformation)
			continue
		}

		parentOrganizationInformation := findParentOrganizationInformation(&organizations, *entity.ParentOrganizationID)
		if parentOrganizationInformation == nil {
			ctx.JSON(http.StatusInternalServerError, errors.New("not found parentOrganizationInformation"))
			return
		}

		if parentOrganizationInformation.SubOrganizations == nil {
			parentOrganizationInformation.SubOrganizations = make([]dtos.OrganizationInformation, 0)
		}

		organizationInformation := factory.NewOrganizationInformationFromEntity(entity)
		parentOrganizationInformation.SubOrganizations = append(parentOrganizationInformation.SubOrganizations, organizationInformation)
	}

	ctx.JSON(http.StatusOK, organizations)
}

func findParentOrganizationInformation(organizations *[]dtos.OrganizationInformation, parentId uint) *dtos.OrganizationInformation {
	for i := 0; i < len(*organizations); i++ {
		if (*organizations)[i].Id == parentId {
			return &(*organizations)[i]
		}
		if (*organizations)[i].SubOrganizations != nil {
			find := findParentOrganizationInformation(&(*organizations)[i].SubOrganizations, parentId)
			if find != nil {
				return find
			}
		}
	}
	return nil
}

func (controller OrganizationController) getOrganization(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	organizationEntity, err := services.OrganizationService{}.GetOrganization(ctx.Request.Context(), uint(organizationId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	organizationRoles := make([]dtos.OrganizationRole, 0)
	for _, role := range organizationEntity.Roles {
		organizationRoles = append(organizationRoles, dtos.OrganizationRole{
			Id:   role.ID,
			Name: role.Name,
		})
	}

	organizationMembers := make([]dtos.OrganizationMember, 0)
	for _, member := range organizationEntity.Members {
		organizationMembers = append(organizationMembers, dtos.OrganizationMember{
			Id:   member.ID,
			Name: member.Name,
		})
	}

	organizationDetails := dtos.OrganizationDetails{
		Id:        organizationEntity.ID,
		Name:      organizationEntity.Name,
		CreatedAt: organizationEntity.CreatedAt,
		Roles:     organizationRoles,
		Members:   organizationMembers,
	}

	ctx.JSON(http.StatusOK, organizationDetails)
}

func (OrganizationController) changeOrganizationName(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.OrganizationService{}.ChangeOrganizationName(ctx.Request.Context(), uint(organizationId), organizationInformation.Name)
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

func (OrganizationController) changePosition(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var requestBody map[string]*uint
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	parentOrganizationId := requestBody["parentOrganizationId"]
	err = services.OrganizationService{}.ChangePosition(ctx.Request.Context(), uint(organizationId), parentOrganizationId)
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

func (OrganizationController) assignRoles(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var organizationAssignRole dtos.OrganizationAssignRole
	if err := ctx.BindJSON(&organizationAssignRole); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.OrganizationService{}.AssignRoles(ctx.Request.Context(), uint(organizationId), organizationAssignRole)
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

func (OrganizationController) assignMembers(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var organizationAssignMember dtos.OrganizationAssignMember
	if err := ctx.BindJSON(&organizationAssignMember); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.OrganizationService{}.AssignMembers(ctx.Request.Context(), uint(organizationId), organizationAssignMember)
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

func (controller OrganizationController) deleteOrganization(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.OrganizationService{}.DeleteOrganization(ctx.Request.Context(), uint(organizationId))
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
