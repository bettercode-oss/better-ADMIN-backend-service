package repository

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/organization/entity"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/pkg/errors"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
}

func (OrganizationRepository) Create(ctx context.Context, entity entity.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (OrganizationRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]entity.OrganizationEntity, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&entity.OrganizationEntity{})

	var entities = make([]entity.OrganizationEntity, 0)

	if err := db.Order("parent_organization_id asc").
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Members").
		Find(&entities).Error; err != nil {
		return entities, errors.Wrap(err, "db error")
	}

	if filters != nil {
		if filters["memberId"] != nil {
			filteredEntities := make([]entity.OrganizationEntity, 0)
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
			filteredEntities := make([]entity.OrganizationEntity, 0)
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

func (OrganizationRepository) FindById(ctx context.Context, id uint) (entity.OrganizationEntity, error) {
	var entity entity.OrganizationEntity

	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Preload("Roles").Preload("Members").First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, errors.Wrap(err, "db error")
	}

	return entity, nil
}

func (OrganizationRepository) Save(ctx context.Context, entity *entity.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Roles").Replace(entity.Roles); err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Model(entity).Association("Members").Replace(entity.Members); err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (OrganizationRepository) Delete(ctx context.Context, entity entity.OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
