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
