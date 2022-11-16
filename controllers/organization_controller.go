package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/factory"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type OrganizationController struct {
}

func (controller OrganizationController) Init(g *echo.Group) {
	g.POST("", controller.CreateOrganization,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.GET("", controller.GetOrganizations,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.GET("/:organizationId", controller.GetOrganization,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.PUT("/:organizationId/name", controller.ChangeOrganizationName,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.PUT("/:organizationId/change-position", controller.ChangePosition,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.PUT("/:organizationId/assign-roles", controller.AssignRoles,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.PUT("/:organizationId/assign-members", controller.AssignMembers,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
	g.DELETE("/:organizationId", controller.DeleteOrganization,
		middlewares.CheckPermission([]string{domain.PermissionManageOrganization}))
}

func (controller OrganizationController) CreateOrganization(ctx echo.Context) error {
	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := organizationInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := services.OrganizationService{}.CreateOrganization(ctx.Request().Context(), organizationInformation)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller OrganizationController) GetOrganizations(ctx echo.Context) error {
	allOfOrganizations, err := services.OrganizationService{}.GetAllOrganizations(ctx.Request().Context(), nil)
	if err != nil {
		return err
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
			return errors.New("not found parentOrganizationInformation")
		}

		if parentOrganizationInformation.SubOrganizations == nil {
			parentOrganizationInformation.SubOrganizations = make([]dtos.OrganizationInformation, 0)
		}

		organizationInformation := factory.NewOrganizationInformationFromEntity(entity)
		parentOrganizationInformation.SubOrganizations = append(parentOrganizationInformation.SubOrganizations, organizationInformation)
	}

	return ctx.JSON(http.StatusOK, organizations)
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

func (OrganizationController) ChangePosition(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var requestBody map[string]*uint
	if err := ctx.Bind(&requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	parentOrganizationId := requestBody["parentOrganizationId"]
	err = services.OrganizationService{}.ChangePosition(ctx.Request().Context(), uint(organizationId), parentOrganizationId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller OrganizationController) DeleteOrganization(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = services.OrganizationService{}.DeleteOrganization(ctx.Request().Context(), uint(organizationId))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (OrganizationController) AssignRoles(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var organizationAssignRole dtos.OrganizationAssignRole
	if err := ctx.Bind(&organizationAssignRole); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := organizationAssignRole.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = services.OrganizationService{}.AssignRoles(ctx.Request().Context(), uint(organizationId), organizationAssignRole)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (OrganizationController) AssignMembers(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var organizationAssignMember dtos.OrganizationAssignMember
	if err := ctx.Bind(&organizationAssignMember); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := organizationAssignMember.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = services.OrganizationService{}.AssignMembers(ctx.Request().Context(), uint(organizationId), organizationAssignMember)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (OrganizationController) ChangeOrganizationName(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := organizationInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = services.OrganizationService{}.ChangeOrganizationName(ctx.Request().Context(), uint(organizationId), organizationInformation.Name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller OrganizationController) GetOrganization(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	organizationEntity, err := services.OrganizationService{}.GetOrganization(ctx.Request().Context(), uint(organizationId))
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusNotFound, nil)
		}
		return err
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

	return ctx.JSON(http.StatusOK, organizationDetails)
}
