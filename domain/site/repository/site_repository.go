package repository

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SiteSettingRepository struct {
}

func (SiteSettingRepository) Save(ctx context.Context, entity entity.SettingEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (SiteSettingRepository) FindByKey(ctx context.Context, key string) (entity.SettingEntity, error) {
	var setting entity.SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&entity.SettingEntity{Key: key}).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return setting, domain.ErrNotFound
		}

		return setting, errors.Wrap(err, "db error")
	}

	return setting, nil
}

func (SiteSettingRepository) FindAll(ctx context.Context) ([]entity.SettingEntity, error) {
	var settings []entity.SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Find(&settings).Error; err != nil {
		return nil, errors.Wrap(err, "db error")
	}

	return settings, nil
}
