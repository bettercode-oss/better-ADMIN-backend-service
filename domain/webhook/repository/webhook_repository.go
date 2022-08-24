package repository

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/webhook/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WebHookRepository struct {
}

func (WebHookRepository) Create(ctx context.Context, entity *entity.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) FindAll(ctx context.Context, pageable dtos.Pageable) ([]entity.WebHookEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&entity.WebHookEntity{})

	var entities = make([]entity.WebHookEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (WebHookRepository) FindById(ctx context.Context, id uint) (entity.WebHookEntity, error) {
	var entity entity.WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.Wrap(err, "db error")
	}

	return entity, nil
}

func (WebHookRepository) Delete(ctx context.Context, entity entity.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) Save(ctx context.Context, entity entity.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) FindLast(ctx context.Context) (entity.WebHookEntity, error) {
	var entity entity.WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.Wrap(err, "db error")
	}

	return entity, nil
}
