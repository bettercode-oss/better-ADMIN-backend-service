package member

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type memberRepository struct {
}

func (r memberRepository) FindBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&MemberEntity{SignId: signId}).
		Preload("Roles.Permissions").Preload(clause.Associations).
		First(&memberEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, errors.Wrap(err, "db error")
	}

	return memberEntity, nil
}

func (memberRepository) FindByDoorayId(ctx context.Context, doorayId string) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&MemberEntity{DoorayId: doorayId}).
		Preload("Roles.Permissions").Preload(clause.Associations).
		First(&memberEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, errors.Wrap(err, "db error")
	}

	return memberEntity, nil
}

func (memberRepository) Create(ctx context.Context, entity *MemberEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

func (memberRepository) FindById(ctx context.Context, id uint) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Preload("Roles.Permissions").Preload(clause.Associations).First(&memberEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, errors.Wrap(err, "db error")
	}

	return memberEntity, nil
}

func (memberRepository) FindAll(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]MemberEntity, int64, error) {
	db := helpers.ContextHelper().GetDB(ctx).Model(&MemberEntity{})

	if filters != nil {
		for key, value := range filters {
			if key == "memberIds" {
				db.Where("id IN ?", value)
			}

			if key == "status" {
				db.Where("status = ?", value)
			}

			if key == "name" {
				db.Where("name LIKE ?", fmt.Sprintf("%%%v%%", value))
			}

			if key == "types" {
				db.Where("type IN ?", value)
			}

			if key == "roleIds" {
				// member_roles ???????????? ???????????? members ????????? ?????? ??? ????????? ??????.
				db.Joins("INNER JOIN member_roles ON member_roles.member_entity_id = members.id").
					Where("member_roles.role_entity_id IN ?", value)
			}
		}
	}

	var entities = make([]MemberEntity, 0)
	var totalCount int64

	if err := db.Count(&totalCount).Scopes(helpers.GormHelper().Pageable(pageable)).
		Preload("Roles.Permissions").Preload(clause.Associations).
		Find(&entities).Error; err != nil {
		return entities, totalCount, errors.Wrap(err, "db error")
	}

	return entities, totalCount, nil
}

func (memberRepository) Save(ctx context.Context, entity *MemberEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Model(entity).Association("Roles").Replace(entity.Roles); err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}

func (memberRepository) FindByGoogleId(ctx context.Context, googleId string) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&MemberEntity{GoogleId: googleId}).
		Preload("Roles.Permissions").Preload(clause.Associations).
		First(&memberEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, errors.Wrap(err, "db error")
	}

	return memberEntity, nil
}

func (memberRepository) Delete(ctx context.Context, entity MemberEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Save(entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	if err := db.Delete(&entity).Error; err != nil {
		return errors.Wrap(err, "db error")
	}

	return nil
}
