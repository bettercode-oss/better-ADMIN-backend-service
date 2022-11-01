package services

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/rbac/entity"
	"better-admin-backend-service/domain/rbac/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/factory"
	"better-admin-backend-service/helpers"
	"context"
)

type RoleBasedAccessControlService struct {
}

func (RoleBasedAccessControlService) CreatePermission(ctx context.Context, permissionInformation dtos.PermissionInformation) error {
	permissionEntity, err := entity.NewPermissionEntity(ctx, permissionInformation)
	if err != nil {
		return err
	}
	return repository.PermissionRepository{}.Create(ctx, &permissionEntity)
}

func (RoleBasedAccessControlService) GetPermissions(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.PermissionEntity, int64, error) {
	return repository.PermissionRepository{}.FindAll(ctx, filters, pageable)
}

func (RoleBasedAccessControlService) GetPermission(ctx context.Context, permissionId uint) (entity.PermissionEntity, error) {
	return repository.PermissionRepository{}.FindById(ctx, permissionId)
}

func (RoleBasedAccessControlService) CreateRole(ctx context.Context, roleInformation dtos.RoleInformation) error {
	roleEntity, err := factory.NewRoleEntity(ctx, roleInformation)
	if err != nil {
		return err
	}

	return repository.RoleRepository{}.Create(ctx, &roleEntity)
}

func (RoleBasedAccessControlService) GetRoles(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.RoleEntity, int64, error) {
	return repository.RoleRepository{}.FindAll(ctx, filters, pageable)
}

func (RoleBasedAccessControlService) UpdatePermission(ctx context.Context, permissionId uint, permissionInformation dtos.PermissionInformation) error {
	repository := repository.PermissionRepository{}
	permissionEntity, err := repository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}

	if permissionEntity.Name != permissionInformation.Name {
		// 변경하려는 이름이 이미 존재하는지 여부 확인
		exists, err := repository.ExistsByName(ctx, permissionInformation.Name)
		if err != nil {
			return err
		}

		if exists == true {
			return domain.ErrDuplicated
		}
	}

	if err := permissionEntity.Update(ctx, permissionInformation); err != nil {
		return err
	}

	return repository.Save(ctx, permissionEntity)
}

func (RoleBasedAccessControlService) DeletePermission(ctx context.Context, permissionId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := repository.PermissionRepository{}

	permissionEntity, err := repository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}

	if err := permissionEntity.Deletable(); err != nil {
		return err
	}

	permissionEntity.UpdatedBy = userClaim.Id

	return repository.Delete(ctx, permissionEntity)
}

func (RoleBasedAccessControlService) DeleteRole(ctx context.Context, roleId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := repository.RoleRepository{}

	roleEntity, err := repository.FindById(ctx, roleId)
	if err != nil {
		return err
	}

	if err := roleEntity.Deletable(); err != nil {
		return err
	}

	roleEntity.UpdatedBy = userClaim.Id

	return repository.Delete(ctx, roleEntity)
}

func (RoleBasedAccessControlService) UpdateRole(ctx context.Context, roleId uint, roleInformation dtos.RoleInformation) error {
	roleRepository := repository.RoleRepository{}
	roleEntity, err := roleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["permissionIds"] = roleInformation.AllowedPermissionIds
	allowedPermissionEntities, _, err := repository.PermissionRepository{}.FindAll(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	if err := roleEntity.Update(ctx, roleInformation, allowedPermissionEntities); err != nil {
		return err
	}

	return roleRepository.Save(ctx, &roleEntity)
}
