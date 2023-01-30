package services

import (
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"better-admin-backend-service/site/domain"
	"better-admin-backend-service/site/repository"
	"context"
	"github.com/mitchellh/mapstructure"
	pkgerrors "github.com/pkg/errors"
)

type SiteService struct {
	siteSettingRepository *repository.SiteSettingRepository
}

func NewSiteService(siteSettingRepository *repository.SiteSettingRepository) *SiteService {
	return &SiteService{
		siteSettingRepository: siteSettingRepository,
	}
}

func (s SiteService) SetSettingWithKey(ctx context.Context, key string, setting interface{}) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		userClaim = &security.UserClaim{
			Id: 0,
		}
	}

	foundSettingEntity, err := s.siteSettingRepository.FindByKey(ctx, key)
	if err != nil {
		if err == errors.ErrNotFound {
			// 설정 값이 없으므로 새로 추가
			newSetting := domain.SettingEntity{
				Key:         key,
				ValueObject: setting,
				CreatedBy:   userClaim.Id,
				UpdatedBy:   userClaim.Id,
			}

			return s.siteSettingRepository.Save(ctx, newSetting)
		}
	}

	foundSettingEntity.ValueObject = setting
	foundSettingEntity.UpdatedBy = userClaim.Id
	return s.siteSettingRepository.Save(ctx, foundSettingEntity)
}

func (s SiteService) GetSettingWithKey(ctx context.Context, key string) (interface{}, error) {
	settingEntity, err := s.siteSettingRepository.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return settingEntity.ValueObject, err
}

func (s SiteService) GetSettings(ctx context.Context) ([]domain.SettingEntity, error) {
	return s.siteSettingRepository.FindAll(ctx)
}

func (s SiteService) GetAppVersion(ctx context.Context) (dtos.AppVersionSetting, error) {
	settingEntity, err := s.siteSettingRepository.FindByKey(ctx, constants.SettingKeyAppVersion)
	if pkgerrors.Is(err, errors.ErrNotFound) {
		newAppVersionSetting := dtos.NewAppVersionSetting()
		if err := s.SetSettingWithKey(ctx, constants.SettingKeyAppVersion, newAppVersionSetting); err != nil {
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

func (s SiteService) IncreaseAppVersion(ctx context.Context) error {
	appVersion, err := s.GetAppVersion(ctx)
	if err != nil {
		return err
	}
	appVersion.Increase()

	return s.SetSettingWithKey(ctx, constants.SettingKeyAppVersion, appVersion)
}
