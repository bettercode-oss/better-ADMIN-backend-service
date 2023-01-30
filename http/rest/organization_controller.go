package rest

import (
	"better-admin-backend-service/app/middlewares"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/organization/factory"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	pkgerrors "github.com/pkg/errors"
	"net/http"
	"strconv"
)

type OrganizationController struct {
	routerGroup         *gin.RouterGroup
	organizationService *services.OrganizationService
}

func NewOrganizationController(
	routerGroup *gin.RouterGroup,
	organizationService *services.OrganizationService) *OrganizationController {

	return &OrganizationController{
		routerGroup:         routerGroup,
		organizationService: organizationService,
	}
}

func (c OrganizationController) MapRoutes() {
	route := c.routerGroup.Group("/organizations")

	route.POST("", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.createOrganization)
	route.GET("", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		middlewares.HttpEtagCache(0),
		c.getOrganizations)
	route.GET("/:organizationId", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		middlewares.HttpEtagCache(0),
		c.getOrganization)
	route.PUT("/:organizationId/name", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.changeOrganizationName)
	route.PUT("/:organizationId/change-position", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.changePosition)
	route.PUT("/:organizationId/assign-roles", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.assignRoles)
	route.PUT("/:organizationId/assign-members", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.assignMembers)
	route.DELETE("/:organizationId", middlewares.PermissionChecker([]string{constants.PermissionManageOrganization}),
		c.deleteOrganization)
}

func (c OrganizationController) createOrganization(ctx *gin.Context) {
	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.organizationService.CreateOrganization(ctx.Request.Context(), organizationInformation)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c OrganizationController) getOrganizations(ctx *gin.Context) {
	allOfOrganizations, err := c.organizationService.GetAllOrganizations(ctx.Request.Context(), nil)
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
			ctx.JSON(http.StatusInternalServerError, pkgerrors.New("not found parentOrganizationInformation"))
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

func (c OrganizationController) getOrganization(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	organizationEntity, err := c.organizationService.GetOrganization(ctx.Request.Context(), uint(organizationId))
	if err != nil {
		if err == errors.ErrNotFound {
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

func (c OrganizationController) changeOrganizationName(ctx *gin.Context) {
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

	err = c.organizationService.ChangeOrganizationName(ctx.Request.Context(), uint(organizationId), organizationInformation.Name)
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

func (c OrganizationController) changePosition(ctx *gin.Context) {
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
	err = c.organizationService.ChangePosition(ctx.Request.Context(), uint(organizationId), parentOrganizationId)
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

func (c OrganizationController) assignRoles(ctx *gin.Context) {
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

	err = c.organizationService.AssignRoles(ctx.Request.Context(), uint(organizationId), organizationAssignRole)
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

func (c OrganizationController) assignMembers(ctx *gin.Context) {
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

	err = c.organizationService.AssignMembers(ctx.Request.Context(), uint(organizationId), organizationAssignMember)
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

func (c OrganizationController) deleteOrganization(ctx *gin.Context) {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.organizationService.DeleteOrganization(ctx.Request.Context(), uint(organizationId))
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
