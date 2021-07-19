package controllers

import (
	"better-admin-backend-service/domain/organization"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type OrganizationController struct {
}

func (controller OrganizationController) Init(g *echo.Group) {
	g.POST("", controller.CreateOrganization, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.GET("", controller.GetOrganizations, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.PUT("/:organizationId/name", controller.ChangeOrganizationName, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.PUT("/:organizationId/change-position", controller.ChangePosition, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.PUT("/:organizationId/assign-roles", controller.AssignRoles, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.PUT("/:organizationId/assign-members", controller.AssignMembers, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
	g.DELETE("/:organizationId", controller.DeleteOrganization, middlewares.CheckPermission([]string{"MANAGE_ORGANIZATION"}))
}

func (controller OrganizationController) CreateOrganization(ctx echo.Context) error {
	var organizationInformation dtos.OrganizationInformation
	if err := ctx.Bind(&organizationInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := organizationInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := organization.OrganizationService{}.CreateOrganization(ctx.Request().Context(), organizationInformation)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller OrganizationController) GetOrganizations(ctx echo.Context) error {
	allOfOrganizations, err := organization.OrganizationService{}.GetAllOrganizations(ctx.Request().Context(), nil)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	organizations := make([]dtos.OrganizationInformation, 0)
	for _, entity := range allOfOrganizations {
		if entity.ParentOrganizationID == nil {
			organizationInformation := dtos.OrganizationInformation{
				Id:   entity.ID,
				Name: entity.Name,
			}

			if entity.Roles != nil && len(entity.Roles) > 0 {
				roles := make([]dtos.OrganizationRole, 0)
				for _, role := range entity.Roles {
					roles = append(roles, dtos.OrganizationRole{
						Id:   role.ID,
						Name: role.Name,
					})
				}
				organizationInformation.OrganizationRoles = roles
			}

			if entity.Members != nil && len(entity.Members) > 0 {
				members := make([]dtos.OrganizationMember, 0)
				for _, member := range entity.Members {
					members = append(members, dtos.OrganizationMember{
						Id:   member.ID,
						Name: member.Name,
					})
				}
				organizationInformation.OrganizationMembers = members
			}
			organizations = append(organizations, organizationInformation)
			continue
		}

		parentOrganizationInformation := findParentOrganizationInformation(&organizations, *entity.ParentOrganizationID)
		if parentOrganizationInformation == nil {
			return ctx.JSON(http.StatusInternalServerError, "not found parentOrganizationInformation")
		}

		if parentOrganizationInformation.SubOrganizations == nil {
			parentOrganizationInformation.SubOrganizations = make([]dtos.OrganizationInformation, 0)
		}

		organizationInformation := dtos.OrganizationInformation{
			Id:   entity.ID,
			Name: entity.Name,
		}
		if entity.Roles != nil && len(entity.Roles) > 0 {
			roles := make([]dtos.OrganizationRole, 0)
			for _, role := range entity.Roles {
				roles = append(roles, dtos.OrganizationRole{
					Id:   role.ID,
					Name: role.Name,
				})
			}
			organizationInformation.OrganizationRoles = roles
		}
		if entity.Members != nil && len(entity.Members) > 0 {
			members := make([]dtos.OrganizationMember, 0)
			for _, member := range entity.Members {
				members = append(members, dtos.OrganizationMember{
					Id:   member.ID,
					Name: member.Name,
				})
			}
			organizationInformation.OrganizationMembers = members
		}
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
	err = organization.OrganizationService{}.ChangePosition(ctx.Request().Context(), uint(organizationId), parentOrganizationId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller OrganizationController) DeleteOrganization(ctx echo.Context) error {
	organizationId, err := strconv.ParseInt(ctx.Param("organizationId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = organization.OrganizationService{}.DeleteOrganization(ctx.Request().Context(), uint(organizationId))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
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

	err = organization.OrganizationService{}.AssignRoles(ctx.Request().Context(), uint(organizationId), organizationAssignRole)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
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

	err = organization.OrganizationService{}.AssignMembers(ctx.Request().Context(), uint(organizationId), organizationAssignMember)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
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

	err = organization.OrganizationService{}.ChangeOrganizationName(ctx.Request().Context(), uint(organizationId), organizationInformation.Name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
