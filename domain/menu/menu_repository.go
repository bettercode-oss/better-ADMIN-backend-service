package menu

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/helpers"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type menuRepository struct {
}

func (menuRepository) Create(ctx context.Context, entity MenuEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return err
	}

	return nil
}

func (menuRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]MenuEntity, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&MenuEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "parentMenuId" {
				if value == nil || reflect.ValueOf(value).IsNil() {
					db.Where("parent_menu_id IS NULL")
				} else {
					db.Where("parent_menu_id = ?", value)
				}
			}
		}
	}

	var entities = make([]MenuEntity, 0)
	if err := db.Order("parent_menu_id asc").Order("sequence asc").Preload(clause.Associations).Find(&entities).Error; err != nil {
		return entities, err
	}

	return entities, nil
}

func (menuRepository) FindById(ctx context.Context, id uint) (MenuEntity, error) {
	var entity MenuEntity

	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, err
	}

	return entity, nil
}

func (menuRepository) Save(ctx context.Context, entity *MenuEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Permissions").Replace(entity.Permissions); err != nil {
		return err
	}

	return db.Save(entity).Error
}

func (menuRepository) Delete(ctx context.Context, entity MenuEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Model(&entity).Association("Permissions").Clear(); err != nil {
		return err
	}

	if err := db.Save(entity).Error; err != nil {
		return err
	}

	return db.Delete(&entity).Error
}
