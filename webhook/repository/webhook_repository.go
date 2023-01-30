package repository

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/webhook/domain"
	"context"
	pkgerrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type WebHookRepository struct {
}

func (WebHookRepository) Create(ctx context.Context, entity *domain.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) FindAll(ctx context.Context, pageable dtos.Pageable) ([]domain.WebHookEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&domain.WebHookEntity{})

	var entities = make([]domain.WebHookEntity, 0)
	var totalCount int64
	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).Find(&entities).Error; err != nil {
		return entities, totalCount, pkgerrors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (WebHookRepository) FindById(ctx context.Context, id uint) (domain.WebHookEntity, error) {
	var entity domain.WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&entity, id).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return entity, errors.ErrNotFound
		}

		return entity, pkgerrors.Wrap(err, "db error")
	}

	return entity, nil
}

func (WebHookRepository) Delete(ctx context.Context, entity domain.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) Save(ctx context.Context, entity domain.WebHookEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (WebHookRepository) FindLast(ctx context.Context) (domain.WebHookEntity, error) {
	var entity domain.WebHookEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Last(&entity).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return entity, errors.ErrNotFound
		}

		return entity, pkgerrors.Wrap(err, "db error")
	}

	return entity, nil
}
