package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AccessControlController struct {
}

func (controller AccessControlController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/access-control")

	route.POST("/permissions", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.createPermission)
	route.GET("/permissions", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		middlewares.HttpEtagCache(0),
		controller.getPermissions)
	route.GET("/permissions/:permissionId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		middlewares.HttpEtagCache(0),
		controller.getPermission)
	route.PUT("/permissions/:permissionId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.updatePermission)
	route.DELETE("/permissions/:permissionId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.deletePermission)
	route.POST("/roles", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.createRole)
	route.GET("/roles", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		middlewares.HttpEtagCache(0),
		controller.getRoles)
	route.GET("/roles/:roleId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		middlewares.HttpEtagCache(0),
		controller.getRole)
	route.PUT("/roles/:roleId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.updateRole)
	route.DELETE("/roles/:roleId", middlewares.PermissionChecker([]string{domain.PermissionManageAccessControl}),
		controller.deleteRole)
}

func (AccessControlController) createPermission(ctx *gin.Context) {
	var permission dtos.PermissionInformation

	if err := ctx.BindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := services.RoleBasedAccessControlService{}.CreatePermission(ctx.Request.Context(), permission)
	if err != nil {
		if err == domain.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AccessControlController) getPermissions(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	filters := map[string]interface{}{}
	if len(ctx.Query("name")) > 0 {
		filters["name"] = ctx.Query("name")
	}

	permissionEntities, totalCount, err := services.RoleBasedAccessControlService{}.GetPermissions(ctx.Request.Context(), filters, pageable)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, pageResult)
}

func (AccessControlController) getPermission(ctx *gin.Context) {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	permissionEntity, err := services.RoleBasedAccessControlService{}.GetPermission(ctx.Request.Context(), uint(permissionId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	permissionDetails := dtos.PermissionDetails{
		Id:          permissionEntity.ID,
		Type:        permissionEntity.Type,
		TypeName:    permissionEntity.GetTypeName(),
		Name:        permissionEntity.Name,
		Description: permissionEntity.Description,
		CreatedAt:   permissionEntity.CreatedAt,
	}

	ctx.JSON(http.StatusOK, permissionDetails)
}

func (AccessControlController) updatePermission(ctx *gin.Context) {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var permission dtos.PermissionInformation
	if err := ctx.BindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RoleBasedAccessControlService{}.UpdatePermission(ctx.Request.Context(), uint(permissionId), permission)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		if err == domain.ErrNonChangeable || err == domain.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AccessControlController) deletePermission(ctx *gin.Context) {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RoleBasedAccessControlService{}.DeletePermission(ctx.Request.Context(), uint(permissionId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == domain.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AccessControlController) createRole(ctx *gin.Context) {
	var role dtos.RoleInformation

	if err := ctx.BindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := services.RoleBasedAccessControlService{}.CreateRole(ctx.Request.Context(), role)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AccessControlController) getRoles(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	filters := map[string]interface{}{}
	if len(ctx.Query("name")) > 0 {
		filters["name"] = ctx.Query("name")
	}

	roleEntities, totalCount, err := services.RoleBasedAccessControlService{}.GetRoles(ctx.Request.Context(), filters, pageable)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, pageResult)
}

func (AccessControlController) getRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	roleEntity, err := services.RoleBasedAccessControlService{}.GetRole(ctx.Request.Context(), uint(roleId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	var allowedPermissions = make([]dtos.AllowedPermission, 0)
	for _, permission := range roleEntity.Permissions {
		allowedPermissions = append(allowedPermissions, dtos.AllowedPermission{
			Id:   permission.ID,
			Name: permission.Name,
		})
	}

	roleDetails := dtos.RoleDetails{
		Id:                 roleEntity.ID,
		Type:               roleEntity.Type,
		TypeName:           roleEntity.GetTypeName(),
		Name:               roleEntity.Name,
		Description:        roleEntity.Description,
		CreatedAt:          roleEntity.CreatedAt,
		AllowedPermissions: allowedPermissions,
	}

	ctx.JSON(http.StatusOK, roleDetails)
}

func (AccessControlController) updateRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var role dtos.RoleInformation

	if err := ctx.BindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RoleBasedAccessControlService{}.UpdateRole(ctx.Request.Context(), uint(roleId), role)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == domain.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AccessControlController) deleteRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RoleBasedAccessControlService{}.DeleteRole(ctx.Request.Context(), uint(roleId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == domain.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
