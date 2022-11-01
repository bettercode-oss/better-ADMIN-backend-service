package repository

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/rbac/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionRepository struct {
}

func (PermissionRepository) Create(ctx context.Context, permissionEntity *entity.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&entity.PermissionEntity{}).Where("name = ?", permissionEntity.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrDuplicated
	}

	if err := db.Create(permissionEntity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

func (PermissionRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.PermissionEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&entity.PermissionEntity{})

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

	var entities = make([]entity.PermissionEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (PermissionRepository) FindById(ctx context.Context, id uint) (entity.PermissionEntity, error) {
	var entity entity.PermissionEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.Wrap(err, "db error")
	}

	return entity, nil
}

func (PermissionRepository) Save(ctx context.Context, entity entity.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	return db.Save(entity).Error
}

func (PermissionRepository) Delete(ctx context.Context, entity entity.PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (PermissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&entity.PermissionEntity{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, "db error")
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

type RoleRepository struct {
}

func (RoleRepository) Create(ctx context.Context, entity *entity.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

func (RoleRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.RoleEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&entity.RoleEntity{})

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

	var entities = make([]entity.RoleEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Preload(clause.Associations).Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (RoleRepository) FindById(ctx context.Context, id uint) (entity.RoleEntity, error) {
	var entity entity.RoleEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Preload(clause.Associations).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.Wrap(err, "db error")
	}

	return entity, nil
}

func (RoleRepository) Delete(ctx context.Context, entity entity.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(&entity).Association("Permissions").Clear(); err != nil {
		return err
	}

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	if err := db.Delete(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (RoleRepository) Save(ctx context.Context, entity *entity.RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
