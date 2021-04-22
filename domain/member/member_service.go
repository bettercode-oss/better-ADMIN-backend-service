package member

import (
	"context"
)

type MemberService struct {
}

func (MemberService) GetMemberBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	return MemberRepository{}.FindBySignId(ctx, signId)
}

func (MemberService) GetMemberByDoorayId(ctx context.Context, doorayId string) (MemberEntity, error) {
	return MemberRepository{}.FindByDoorayId(ctx, doorayId)
}

func (MemberService) CreateMember(ctx context.Context, entity *MemberEntity) error {
	return MemberRepository{}.CreateMember(ctx, entity)
}

func (MemberService) GetMemberById(ctx context.Context, id uint) (MemberEntity, error) {
	return MemberRepository{}.FindById(ctx, id)
}
