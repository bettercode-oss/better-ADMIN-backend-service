package member

import (
	"context"
)

type MemberService struct {
}

func (service MemberService) GetMemberBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	return MemberRepository{}.FindBySignId(ctx, signId)
}
