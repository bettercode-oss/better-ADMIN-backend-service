package services

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/domain/site/repository"
	"better-admin-backend-service/helpers"
	"context"
)

type SiteService struct {
}

func (SiteService) SetSettingWithKey(ctx context.Context, key string, setting interface{}) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := repository.SiteSettingRepository{}
	foundSettingEntity, err := repository.FindByKey(ctx, key)
	if err != nil {
		if err == domain.ErrNotFound {
			// 설정 값이 없으므로 새로 추가
			newSetting := entity.SettingEntity{
				Key:         key,
				ValueObject: setting,
				CreatedBy:   userClaim.Id,
				UpdatedBy:   userClaim.Id,
			}

			return repository.Save(ctx, newSetting)
		}
	}

	foundSettingEntity.ValueObject = setting
	foundSettingEntity.UpdatedBy = userClaim.Id
	return repository.Save(ctx, foundSettingEntity)
}

func (SiteService) GetSettingWithKey(ctx context.Context, key string) (interface{}, error) {
	settingEntity, err := repository.SiteSettingRepository{}.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return settingEntity.ValueObject, err
}

func (SiteService) GetSettings(ctx context.Context) ([]entity.SettingEntity, error) {
	return repository.SiteSettingRepository{}.FindAll(ctx)
}
