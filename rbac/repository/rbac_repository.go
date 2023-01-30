package repository

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/rbac/domain"
	"context"
	"fmt"
	pkgerrors "github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionRepository struct {
}

func (PermissionRepository) Create(ctx context.Context, permissionEntity *domain.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&domain.PermissionEntity{}).Where("name = ?", permissionEntity.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.ErrDuplicated
	}

	if err := db.Create(permissionEntity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}
	return nil
}

func (PermissionRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]domain.PermissionEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&domain.PermissionEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "permissionIds" {
				db.Where("id IN ?", value)
			}

			if key == "name" {
				db.Where("name LIKE ?", fmt.Sprintf("%%%v%%", value))
			}
		}
	}

	var entities = make([]domain.PermissionEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, pkgerrors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (PermissionRepository) FindById(ctx context.Context, id uint) (domain.PermissionEntity, error) {
	var entity domain.PermissionEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return entity, errors.ErrNotFound
		}

		return entity, pkgerrors.Wrap(err, "db error")
	}

	return entity, nil
}

func (PermissionRepository) Save(ctx context.Context, entity domain.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	return db.Save(entity).Error
}

func (PermissionRepository) Delete(ctx context.Context, entity domain.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (PermissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&domain.PermissionEntity{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, pkgerrors.Wrap(err, "db error")
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

type RoleRepository struct {
}

func (RoleRepository) Create(ctx context.Context, entity *domain.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}
	return nil
}

func (RoleRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]domain.RoleEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&domain.RoleEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "roleIds" {
				db.Where("id IN ?", value)
			}

			if key == "name" {
				db.Where("name LIKE ?", fmt.Sprintf("%%%v%%", value))
			}
		}
	}

	var entities = make([]domain.RoleEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Preload(clause.Associations).Find(&entities).Error; err != nil {
		return entities, totalCount, pkgerrors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (RoleRepository) FindById(ctx context.Context, id uint) (domain.RoleEntity, error) {
	var entity domain.RoleEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Preload(clause.Associations).First(&entity, id).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return entity, errors.ErrNotFound
		}

		return entity, pkgerrors.Wrap(err, "db error")
	}

	return entity, nil
}

func (RoleRepository) Delete(ctx context.Context, entity domain.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(&entity).Association("Permissions").Clear(); err != nil {
		return err
	}

	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}
	if err := db.Delete(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (RoleRepository) Save(ctx context.Context, entity *domain.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}
