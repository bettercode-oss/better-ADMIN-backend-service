package services

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/rbac/domain"
	"better-admin-backend-service/rbac/factory"
	"better-admin-backend-service/rbac/repository"
	"context"
)

type RoleBasedAccessControlService struct {
	permissionRepository *repository.PermissionRepository
	roleRepository       *repository.RoleRepository
}

func NewRoleBasedAccessControlService(
	permissionRepository *repository.PermissionRepository,
	roleRepository *repository.RoleRepository) *RoleBasedAccessControlService {

	return &RoleBasedAccessControlService{
		permissionRepository: permissionRepository,
		roleRepository:       roleRepository,
	}
}

func (s RoleBasedAccessControlService) CreatePermission(ctx context.Context, permissionInformation dtos.PermissionInformation) error {
	permissionEntity, err := domain.NewPermissionEntity(ctx, permissionInformation)
	if err != nil {
		return err
	}
	return s.permissionRepository.Create(ctx, &permissionEntity)
}

func (s RoleBasedAccessControlService) GetPermissions(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]domain.PermissionEntity, int64, error) {
	return s.permissionRepository.FindAll(ctx, filters, pageable)
}

func (s RoleBasedAccessControlService) GetPermission(ctx context.Context, permissionId uint) (domain.PermissionEntity, error) {
	return s.permissionRepository.FindById(ctx, permissionId)
}

func (s RoleBasedAccessControlService) CreateRole(ctx context.Context, roleInformation dtos.RoleInformation) error {
	roleEntity, err := factory.NewRoleEntity(ctx, roleInformation, s.permissionRepository)
	if err != nil {
		return err
	}

	return s.roleRepository.Create(ctx, &roleEntity)
}

func (s RoleBasedAccessControlService) GetRoles(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]domain.RoleEntity, int64, error) {
	return s.roleRepository.FindAll(ctx, filters, pageable)
}

func (s RoleBasedAccessControlService) UpdatePermission(ctx context.Context, permissionId uint, permissionInformation dtos.PermissionInformation) error {
	permissionEntity, err := s.permissionRepository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}

	if permissionEntity.Name != permissionInformation.Name {
		// 변경하려는 이름이 이미 존재하는지 여부 확인
		exists, err := s.permissionRepository.ExistsByName(ctx, permissionInformation.Name)
		if err != nil {
			return err
		}

		if exists == true {
			return errors.ErrDuplicated
		}
	}

	if err := permissionEntity.Update(ctx, permissionInformation); err != nil {
		return err
	}

	return s.permissionRepository.Save(ctx, permissionEntity)
}

func (s RoleBasedAccessControlService) DeletePermission(ctx context.Context, permissionId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	permissionEntity, err := s.permissionRepository.FindById(ctx, permissionId)
	if err != nil {
		return err
	}

	if err := permissionEntity.Deletable(); err != nil {
		return err
	}

	permissionEntity.UpdatedBy = userClaim.Id

	return s.permissionRepository.Delete(ctx, permissionEntity)
}

func (s RoleBasedAccessControlService) DeleteRole(ctx context.Context, roleId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	roleEntity, err := s.roleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}

	if err := roleEntity.Deletable(); err != nil {
		return err
	}

	roleEntity.UpdatedBy = userClaim.Id

	return s.roleRepository.Delete(ctx, roleEntity)
}

func (s RoleBasedAccessControlService) UpdateRole(ctx context.Context, roleId uint, roleInformation dtos.RoleInformation) error {
	roleEntity, err := s.roleRepository.FindById(ctx, roleId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["permissionIds"] = roleInformation.AllowedPermissionIds
	allowedPermissionEntities, _, err := s.permissionRepository.FindAll(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	if err := roleEntity.Update(ctx, roleInformation, allowedPermissionEntities); err != nil {
		return err
	}

	return s.roleRepository.Save(ctx, &roleEntity)
}

func (s RoleBasedAccessControlService) GetRole(ctx context.Context, roleId uint) (domain.RoleEntity, error) {
	return s.roleRepository.FindById(ctx, roleId)
}
