package site

import (
	"better-admin-backend-service/domain"
	"context"
)

type SiteService struct {
}

func (SiteService) SetSettingWithKey(ctx context.Context, key string, setting interface{}) error {
	repository := SiteSettingRepository{}

	foundSettingEntity, err := repository.FindByKey(ctx, key)
	if err != nil {
		if err == domain.ErrNotFound {
			// 설정 값이 없으므로 새로 추가
			newSetting := SettingEntity{
				Key:         key,
				ValueObject: setting,
			}

			return repository.Save(ctx, newSetting)
		}
	}

	foundSettingEntity.ValueObject = setting
	return repository.Save(ctx, foundSettingEntity)
}

func (SiteService) GetSettingWithKey(ctx context.Context, key string) (interface{}, error) {
	settingEntity, err := SiteSettingRepository{}.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return settingEntity.ValueObject, err
}

func (SiteService) GetSettings(ctx context.Context) ([]SettingEntity, error) {
	return SiteSettingRepository{}.FindAll(ctx)
}
