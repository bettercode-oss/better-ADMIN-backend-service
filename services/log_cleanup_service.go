package services

import (
	"better-admin-backend-service/domain/logging/repository"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/dtos"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	logCleanupServiceOnce     sync.Once
	logCleanupServiceInstance *logCleanupService
)

func LogCleanupService() *logCleanupService {
	logCleanupServiceOnce.Do(func() {
		logCleanupServiceInstance = &logCleanupService{}
	})

	return logCleanupServiceInstance
}

type logCleanupService struct {
	cleanupDate *time.Time
}

func (service *logCleanupService) cleanupDaily(ctx context.Context) {
	// 성능 상 매번 하지 않고 하루에 한번만 실행한다.
	now := time.Now()
	if service.cleanupDate == nil || service.cleanupDate.AddDate(0, 0, 1).Before(now) {
		if err := service.deleteLogs(ctx); err != nil {
			log.Error("cleanup logs error : ", err)
		} else {
			service.cleanupDate = &now
		}
	}
}

func (logCleanupService) deleteLogs(ctx context.Context) error {
	setting, err := SiteService{}.GetSettingWithKey(ctx, entity.SettingKeyMemberAccessLog)
	if err != nil {
		return err
	}

	var memberAccessLogSetting dtos.MemberAccessLogSetting
	if err := mapstructure.Decode(setting, &memberAccessLogSetting); err != nil {
		return errors.Wrap(err, "map to struct decode error")
	}

	retentionDays := memberAccessLogSetting.RetentionDays
	beforeDateOfDeletion := time.Now().AddDate(0, 0, -int(retentionDays))

	return repository.MemberAccessLogRepository{}.DeleteBeforeDate(ctx, beforeDateOfDeletion)
}
