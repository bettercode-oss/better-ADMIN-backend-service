package repository

import (
	"better-admin-backend-service/domain/logging/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/pkg/errors"
	"time"
)

type MemberAccessLogRepository struct {
}

func (MemberAccessLogRepository) Create(ctx context.Context, entity entity.MemberAccessLogEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

func (MemberAccessLogRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.MemberAccessLogEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&entity.MemberAccessLogEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "memberId" {
				db.Where("member_id = ?", value)
			}
		}
	}

	var entities = make([]entity.MemberAccessLogEntity, 0)
	var totalCount int64

	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Order("id desc").Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (MemberAccessLogRepository) DeleteBeforeDate(ctx context.Context, date time.Time) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Where("created_at < ?", date).Delete(&entity.MemberAccessLogEntity{}).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
