package factory

import (
	"better-admin-backend-service/domain/rbac/entity"
	"better-admin-backend-service/domain/rbac/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
)

func NewRoleEntity(ctx context.Context, information dtos.RoleInformation) (entity.RoleEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return entity.RoleEntity{}, err
	}

	role := entity.RoleEntity{
		Type:        entity.UserDefineTypeKey,
		Name:        information.Name,
		Description: information.Description,
		CreatedBy:   userClaim.Id,
		UpdatedBy:   userClaim.Id,
	}
	filters := map[string]interface{}{}
	filters["permissionIds"] = information.AllowedPermissionIds

	permissionEntities, _, err := repository.PermissionRepository{}.FindAll(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return role, err
	}

	role.Permissions = permissionEntities
	return role, nil
}
