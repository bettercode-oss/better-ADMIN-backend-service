package services

import (
	"better-admin-backend-service/domain/logging/entity"
	"better-admin-backend-service/domain/logging/repository"
	"better-admin-backend-service/dtos"
	"context"
)

type MemberAccessLogService struct {
}

func (MemberAccessLogService) LogMemberAccess(ctx context.Context, accessLog dtos.MemberAccessLog) error {
	LogCleanupService().cleanupDaily(ctx)
	entity, err := entity.NewMemberAccessLogEntity(ctx, accessLog)
	if err != nil {
		return err
	}

	return repository.MemberAccessLogRepository{}.Create(ctx, entity)
}

func (MemberAccessLogService) GetMemberAccessLogs(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.MemberAccessLogEntity, int64, error) {
	return repository.MemberAccessLogRepository{}.FindAll(ctx, filters, pageable)
}

func (service MemberAccessLogService) cleanUpLog() {

}
