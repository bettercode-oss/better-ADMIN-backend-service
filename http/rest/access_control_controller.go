package rest

import (
	"better-admin-backend-service/app/middlewares"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/services"
	etag "github.com/bettercode-oss/gin-middleware-etag"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AccessControlController struct {
	routerGroup                   *gin.RouterGroup
	roleBasedAccessControlService *services.RoleBasedAccessControlService
}

func NewAccessControlController(rg *gin.RouterGroup,
	roleBasedAccessControlService *services.RoleBasedAccessControlService) *AccessControlController {
	return &AccessControlController{
		routerGroup:                   rg,
		roleBasedAccessControlService: roleBasedAccessControlService,
	}
}

func (c AccessControlController) MapRoutes() {
	route := c.routerGroup.Group("/access-control")

	route.POST("/permissions", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.createPermission)
	route.GET("/permissions", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		etag.HttpEtagCache(0),
		c.getPermissions)
	route.GET("/permissions/:permissionId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		etag.HttpEtagCache(0),
		c.getPermission)
	route.PUT("/permissions/:permissionId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.updatePermission)
	route.DELETE("/permissions/:permissionId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.deletePermission)
	route.POST("/roles", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.createRole)
	route.GET("/roles", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		etag.HttpEtagCache(0),
		c.getRoles)
	route.GET("/roles/:roleId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		etag.HttpEtagCache(0),
		c.getRole)
	route.PUT("/roles/:roleId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.updateRole)
	route.DELETE("/roles/:roleId", middlewares.PermissionChecker([]string{constants.PermissionManageAccessControl}),
		c.deleteRole)
}

func (c AccessControlController) createPermission(ctx *gin.Context) {
	var permission dtos.PermissionInformation

	if err := ctx.BindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.roleBasedAccessControlService.CreatePermission(ctx.Request.Context(), permission)
	if err != nil {
		if err == errors.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c AccessControlController) getPermissions(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	filters := map[string]interface{}{}
	if len(ctx.Query("name")) > 0 {
		filters["name"] = ctx.Query("name")
	}

	permissionEntities, totalCount, err := c.roleBasedAccessControlService.GetPermissions(ctx.Request.Context(), filters, pageable)
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

func (c AccessControlController) getPermission(ctx *gin.Context) {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	permissionEntity, err := c.roleBasedAccessControlService.GetPermission(ctx.Request.Context(), uint(permissionId))
	if err != nil {
		if err == errors.ErrNotFound {
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

func (c AccessControlController) updatePermission(ctx *gin.Context) {
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

	err = c.roleBasedAccessControlService.UpdatePermission(ctx.Request.Context(), uint(permissionId), permission)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		if err == errors.ErrNonChangeable || err == errors.ErrDuplicated {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c AccessControlController) deletePermission(ctx *gin.Context) {
	permissionId, err := strconv.ParseInt(ctx.Param("permissionId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.roleBasedAccessControlService.DeletePermission(ctx.Request.Context(), uint(permissionId))
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == errors.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c AccessControlController) createRole(ctx *gin.Context) {
	var role dtos.RoleInformation

	if err := ctx.BindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.roleBasedAccessControlService.CreateRole(ctx.Request.Context(), role)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c AccessControlController) getRoles(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	filters := map[string]interface{}{}
	if len(ctx.Query("name")) > 0 {
		filters["name"] = ctx.Query("name")
	}

	roleEntities, totalCount, err := c.roleBasedAccessControlService.GetRoles(ctx.Request.Context(), filters, pageable)
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

func (c AccessControlController) getRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	roleEntity, err := c.roleBasedAccessControlService.GetRole(ctx.Request.Context(), uint(roleId))
	if err != nil {
		if err == errors.ErrNotFound {
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

func (c AccessControlController) updateRole(ctx *gin.Context) {
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

	err = c.roleBasedAccessControlService.UpdateRole(ctx.Request.Context(), uint(roleId), role)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == errors.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c AccessControlController) deleteRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.roleBasedAccessControlService.DeleteRole(ctx.Request.Context(), uint(roleId))
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		if err == errors.ErrNonChangeable {
			ctx.JSON(http.StatusBadRequest, dtos.ErrorMessage{Message: err.Error()})
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
