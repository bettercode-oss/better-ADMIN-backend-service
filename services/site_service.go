package services

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/domain/site/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type SiteService struct {
}

func (SiteService) SetSettingWithKey(ctx context.Context, key string, setting interface{}) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		userClaim = &security.UserClaim{
			Id: 0,
		}
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

func (service SiteService) GetAppVersion(ctx context.Context) (dtos.AppVersionSetting, error) {
	settingEntity, err := repository.SiteSettingRepository{}.FindByKey(ctx, entity.SettingKeyAppVersion)
	if errors.Is(err, domain.ErrNotFound) {
		newAppVersionSetting := dtos.NewAppVersionSetting()
		if err := service.SetSettingWithKey(ctx, entity.SettingKeyAppVersion, newAppVersionSetting); err != nil {
			return dtos.AppVersionSetting{}, err
		}

		return newAppVersionSetting, nil
	}

	var appVersion dtos.AppVersionSetting
	if err = mapstructure.Decode(settingEntity.ValueObject, &appVersion); err != nil {
		return dtos.AppVersionSetting{}, err
	}

	return appVersion, nil
}

func (service SiteService) IncreaseAppVersion(ctx context.Context) error {
	appVersion, err := service.GetAppVersion(ctx)
	if err != nil {
		return err
	}
	appVersion.Increase()

	return service.SetSettingWithKey(ctx, entity.SettingKeyAppVersion, appVersion)
}
