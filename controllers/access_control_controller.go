package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type AccessControlController struct {
}

func (controller AccessControlController) Init(g *echo.Group) {
	g.POST("/permissions", controller.CreatePermission,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.GET("/permissions", controller.GetPermissions,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.PUT("/permissions/:permissionId", controller.UpdatePermission,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.DELETE("/permissions/:permissionId", controller.DeletePermission,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.POST("/roles", controller.CreateRole,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.GET("/roles", controller.GetRoles,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.PUT("/roles/:roleId", controller.UpdateRole,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
	g.DELETE("/roles/:roleId", controller.DeleteRole,
		middlewares.CheckPermission([]string{domain.PermissionManageAccessControl}))
}

func (AccessControlController) CreatePermission(ctx echo.Context) error {
	var permission dtos.PermissionInformation

	if err := ctx.Bind(&permission); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := permission.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := rbac.RoleBasedAccessControlService{}.CreatePermission(ctx.Request().Context(), permission)
	if err != nil {
		if err == domain.ErrDuplicated {
			return ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
		}

		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (AccessControlController) GetPermissions(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)

	permissionEntities, totalCount, err := rbac.RoleBasedAccessControlService{}.GetPermissions(ctx.Request().Context(), nil, pageable)
	if err != nil {
		return err
	}

	var permissions = make([]dtos.PermissionInformation, 0)
	for _, entity := range permissionEntities {
		permissions = append(permissions, dtos.PermissionInformation{
			Id:          entity.ID,
			Type:        entity.Type,
			TypeName:    entity.GetTypeName(),
			Name:        entity.Name,
			Description: entity.Description,
		})
	}

	pageResult := dtos.PageResult{
		Result:     permissions,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, pageResult)
}

func (AccessControlController) UpdatePermission(ctx echo.Context) error {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var permission dtos.PermissionInformation
	if err := ctx.Bind(&permission); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := permission.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = rbac.RoleBasedAccessControlService{}.UpdatePermission(ctx.Request().Context(), uint(permissionId), permission)
	if err != nil {
		if err == domain.ErrNonChangeable || err == domain.ErrDuplicated {
			return ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
		}
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (AccessControlController) DeletePermission(ctx echo.Context) error {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = rbac.RoleBasedAccessControlService{}.DeletePermission(ctx.Request().Context(), uint(permissionId))
	if err != nil {
		if err == domain.ErrNonChangeable {
			return ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
		}
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (AccessControlController) CreateRole(ctx echo.Context) error {
	var role dtos.RoleInformation

	if err := ctx.Bind(&role); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := role.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := rbac.RoleBasedAccessControlService{}.CreateRole(ctx.Request().Context(), role)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (AccessControlController) GetRoles(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)

	roleEntities, totalCount, err := rbac.RoleBasedAccessControlService{}.GetRoles(ctx.Request().Context(), nil, pageable)
	if err != nil {
		return err
	}

	var roleSummaries = make([]dtos.RoleSummary, 0)
	for _, role := range roleEntities {
		var allowedPermissions = make([]dtos.AllowedPermission, 0)
		for _, permission := range role.Permissions {
			allowedPermissions = append(allowedPermissions, dtos.AllowedPermission{
				Id:   permission.ID,
				Name: permission.Name,
			})
		}

		roleSummaries = append(roleSummaries, dtos.RoleSummary{
			Id:                role.ID,
			Type:              role.Type,
			TypeName:          role.GetTypeName(),
			Name:              role.Name,
			Description:       role.Description,
			AllowedPermission: allowedPermissions,
		})
	}

	pageResult := dtos.PageResult{
		Result:     roleSummaries,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, pageResult)
}

func (AccessControlController) DeleteRole(ctx echo.Context) error {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = rbac.RoleBasedAccessControlService{}.DeleteRole(ctx.Request().Context(), uint(roleId))
	if err != nil {
		if err == domain.ErrNonChangeable {
			return ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
		}
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (AccessControlController) UpdateRole(ctx echo.Context) error {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var role dtos.RoleInformation

	if err := ctx.Bind(&role); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := role.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = rbac.RoleBasedAccessControlService{}.UpdateRole(ctx.Request().Context(), uint(roleId), role)
	if err != nil {
		if err == domain.ErrNonChangeable {
			return ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
		}
		return err
	}

	return ctx.JSON(http.StatusOK, nil)

}
