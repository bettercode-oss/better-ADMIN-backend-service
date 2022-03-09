package site

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type siteSettingRepository struct {
}

func (siteSettingRepository) Save(ctx context.Context, entity SettingEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(&entity).Error; err != nil {
		return errors.New(err)
	}

	return nil
}

func (siteSettingRepository) FindByKey(ctx context.Context, key string) (SettingEntity, error) {
	var setting SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&SettingEntity{Key: key}).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return setting, domain.ErrNotFound
		}

		return setting, errors.New(err)
	}

	return setting, nil
}

func (siteSettingRepository) FindAll(ctx context.Context) ([]SettingEntity, error) {
	var settings []SettingEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Find(&settings).Error; err != nil {
		return nil, errors.New(err)
	}

	return settings, nil
}
