package member

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/helpers"
	"context"
	"errors"
	"gorm.io/gorm"
)

type MemberRepository struct {
}

func (r MemberRepository) FindBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&MemberEntity{SignId: signId}).First(&memberEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, err
	}

	return memberEntity, nil
}

func (MemberRepository) FindByDoorayId(ctx context.Context, doorayId string) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.Where(&MemberEntity{DoorayId: doorayId}).First(&memberEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, err
	}

	return memberEntity, nil
}

func (MemberRepository) CreateMember(ctx context.Context, entity *MemberEntity) error {
	db := helpers.ContextHelper().GetDB(ctx)
	if err := db.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (MemberRepository) FindById(ctx context.Context, id uint) (MemberEntity, error) {
	var memberEntity MemberEntity

	db := helpers.ContextHelper().GetDB(ctx)

	if err := db.First(&memberEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return memberEntity, domain.ErrNotFound
		}

		return memberEntity, err
	}

	return memberEntity, nil
}
