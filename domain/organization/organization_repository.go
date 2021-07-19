package organization

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/helpers"
	"context"
	"errors"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
)

type organizationRepository struct {
}

func (organizationRepository) Create(ctx context.Context, entity OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(&entity).Error; err != nil {
		return err
	}

	return nil
}

func (organizationRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]OrganizationEntity, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&OrganizationEntity{})

	var entities = make([]OrganizationEntity, 0)

	if err := db.Order("parent_organization_id asc").
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Members").
		Find(&entities).Error; err != nil {
		return entities, err
	}

	if filters != nil {
		if filters["memberId"] != nil {
			filteredEntities := make([]OrganizationEntity, 0)
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
			filteredEntities := make([]OrganizationEntity, 0)
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

func (organizationRepository) FindById(ctx context.Context, id uint) (OrganizationEntity, error) {
	var entity OrganizationEntity

	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Preload("Roles").Preload("Members").First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, domain.ErrNotFound
		}

		return entity, err
	}

	return entity, nil
}

func (organizationRepository) Save(ctx context.Context, entity *OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Roles").Replace(entity.Roles); err != nil {
		return err
	}

	if err := db.Model(entity).Association("Members").Replace(entity.Members); err != nil {
		return err
	}

	return db.Save(entity).Error
}

func (organizationRepository) Delete(ctx context.Context, entity OrganizationEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	return db.Delete(&entity).Error
}
