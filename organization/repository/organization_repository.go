package repository

import (
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/organization/domain"
	"context"
	pkgerrors "github.com/pkg/errors"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
}

func (OrganizationRepository) Create(ctx context.Context, entity domain.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (OrganizationRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]domain.OrganizationEntity, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&domain.OrganizationEntity{})

	var entities = make([]domain.OrganizationEntity, 0)

	if err := db.Order("parent_organization_id asc").
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Members").
		Find(&entities).Error; err != nil {
		return entities, pkgerrors.Wrap(err, "db error")
	}

	if filters != nil {
		if filters["memberId"] != nil {
			filteredEntities := make([]domain.OrganizationEntity, 0)
			for _, entity := range entities {
				for _, member := range entity.Members {
					if member.ID == filters["memberId"].(uint) {
						filteredEntities = append(filteredEntities, entity)
						break
					}
				}
			}

			return filteredEntities, nil
		}

		if filters["memberIds"] != nil {
			filteredEntities := make([]domain.OrganizationEntity, 0)
			stream := koazee.StreamOf(filters["memberIds"].([]uint))
			for _, entity := range entities {
				for _, member := range entity.Members {
					contains, _ := stream.Contains(member.ID)
					if contains {
						filteredEntities = append(filteredEntities, entity)
						break
					}
				}
			}

			return filteredEntities, nil
		}

	}

	return entities, nil
}

func (OrganizationRepository) FindById(ctx context.Context, id uint) (domain.OrganizationEntity, error) {
	var entity domain.OrganizationEntity

	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Preload("Roles").Preload("Members").First(&entity, id).Error; err != nil {
		if pkgerrors.Is(err, gorm.ErrRecordNotFound) {
			return entity, errors.ErrNotFound
		}

		return entity, pkgerrors.Wrap(err, "db error")
	}

	return entity, nil
}

func (OrganizationRepository) Save(ctx context.Context, entity *domain.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Roles").Replace(entity.Roles); err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Model(entity).Association("Members").Replace(entity.Members); err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}

func (OrganizationRepository) Delete(ctx context.Context, entity domain.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return pkgerrors.Wrap(err, "db error")
	}

	return nil
}
