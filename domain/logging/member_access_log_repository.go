package logging

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
)

type memberAccessLogRepository struct {
}

func (memberAccessLogRepository) Create(ctx context.Context, entity MemberAccessLogEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return err
	}
	return nil
}

func (memberAccessLogRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]MemberAccessLogEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&MemberAccessLogEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "memberId" {
				db.Where("member_id = ?", value)
			}
		}
	}

	var entities = make([]MemberAccessLogEntity, 0)
	var totalCount int64

	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Order("id desc").Find(&entities).Error; err != nil {
		return entities, totalCount, err
	}

	return entities, totalCount, nil
}
