package factory

import (
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/rbac/domain"
	"better-admin-backend-service/rbac/repository"
	"context"
)

func NewRoleEntity(ctx context.Context, information dtos.RoleInformation, permissionRepository *repository.PermissionRepository) (domain.RoleEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return domain.RoleEntity{}, err
	}

	role := domain.RoleEntity{
		Type:        constants.UserDefineTypeKey,
		Name:        information.Name,
		Description: information.Description,
		CreatedBy:   userClaim.Id,
		UpdatedBy:   userClaim.Id,
	}
	filters := map[string]interface{}{}
	filters["permissionIds"] = information.AllowedPermissionIds

	permissionEntities, _, err := permissionRepository.FindAll(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return role, err
	}

	role.Permissions = permissionEntities
	return role, nil
}
