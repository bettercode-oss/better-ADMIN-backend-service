package rbac

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type permissionRepository struct {
}

func (permissionRepository) Create(ctx context.Context, entity *PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&PermissionEntity{}).Where("name = ?", entity.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrDuplicated
	}

	if err := db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (permissionRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]PermissionEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&PermissionEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "permissionIds" {
				db.Where("id IN ?", value)
			}
		}
	}

	var entities = make([]PermissionEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, err
	}

	return entities, totalCount, nil
}

func (permissionRepository) FindById(ctx context.Context, id uint) (PermissionEntity, error) {
	var entity PermissionEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, err
	}

	return entity, nil
}

func (permissionRepository) Save(ctx context.Context, entity PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	return db.Save(entity).Error
}

func (permissionRepository) Delete(ctx context.Context, entity PermissionEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return err
	}
	return db.Delete(&entity).Error
}

func (permissionRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	db := helpers.ContextHelper().GetDB(ctx)

	var count int64
	if err := db.Model(&PermissionEntity{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

type roleRepository struct {
}

func (roleRepository) Create(ctx context.Context, entity *RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (roleRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]RoleEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&RoleEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "roleIds" {
				db.Where("id IN ?", value)
			}
		}
	}

	var entities = make([]RoleEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Preload(clause.Associations).Find(&entities).Error; err != nil {
		return entities, totalCount, err
	}

	return entities, totalCount, nil
}

func (roleRepository) FindById(ctx context.Context, id uint) (RoleEntity, error) {
	var entity RoleEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Preload(clause.Associations).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, err
	}

	return entity, nil
}

func (roleRepository) Delete(ctx context.Context, entity RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return err
	}
	return db.Delete(&entity).Error
}

func (roleRepository) Save(ctx context.Context, entity *RoleEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
		return err
	}

	return db.Save(entity).Error
}
