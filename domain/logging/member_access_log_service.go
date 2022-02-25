package logging

import (
	"better-admin-backend-service/dtos"
	"context"
)

type MemberAccessLogService struct {
}

func (MemberAccessLogService) LogMemberAccess(ctx context.Context, accessLog dtos.MemberAccessLog) error {
	entity, err := NewMemberAccessLogEntity(ctx, accessLog)
	if err != nil {
		return err
	}

	return memberAccessLogRepository{}.Create(ctx, entity)
}

func (MemberAccessLogService) GetMemberAccessLogs(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]MemberAccessLogEntity, int64, error) {
	return memberAccessLogRepository{}.FindAll(ctx, filters, pageable)
}
