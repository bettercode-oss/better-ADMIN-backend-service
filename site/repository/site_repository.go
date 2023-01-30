package repository

import (
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/site/domain"
	"context"
	pkgerrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type SiteSettingRepository struct {
}

func (SiteSettingRepository) Save(ctx context.Context, entity domain.SettingEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (SiteSettingRepository) FindByKey(ctx context.Context, key string) (domain.SettingEntity, error) {
	var setting domain.SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&domain.SettingEntity{Key: key}).First(&setting).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return setting, errors.ErrNotFound
		}

		return setting, pkgerrors.Wrap(err, "db error")
	}

	return setting, nil
}

func (SiteSettingRepository) FindAll(ctx context.Context) ([]domain.SettingEntity, error) {
	var settings []domain.SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Find(&settings).Error; err != nil {
		return nil, pkgerrors.Wrap(err, "db error")
	}

	return settings, nil
}
