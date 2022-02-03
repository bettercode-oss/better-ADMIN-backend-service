package member

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
)

type MemberService struct {
}

func (MemberService) GetMemberBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	return memberRepository{}.FindBySignId(ctx, signId)
}

func (MemberService) GetMemberByDoorayId(ctx context.Context, doorayId string) (MemberEntity, error) {
	return memberRepository{}.FindByDoorayId(ctx, doorayId)
}

func (MemberService) CreateMember(ctx context.Context, entity *MemberEntity) error {
	return memberRepository{}.Create(ctx, entity)
}

func (MemberService) GetMemberById(ctx context.Context, id uint) (MemberEntity, error) {
	return memberRepository{}.FindById(ctx, id)
}

func (MemberService) GetMembers(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]MemberEntity, int64, error) {
	return memberRepository{}.FindAll(ctx, filters, pageable)
}

func (MemberService) AssignRole(ctx context.Context, memberId uint, assignRole dtos.MemberAssignRole) error {
	repository := memberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	err = memberEntity.AssignRole(ctx, assignRole)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &memberEntity)
}

func (MemberService) GetMember(ctx context.Context, memberId uint) (MemberEntity, error) {
	return memberRepository{}.FindById(ctx, memberId)
}

func (MemberService) SignUpMember(ctx context.Context, signUp dtos.MemberSignUp) error {
	repository := memberRepository{}
	_, err := repository.FindBySignId(ctx, signUp.SignId)
	if err != nil {
		if err == domain.ErrNotFound {
			// signId 가 중복이 없을 때만 가입
			newMember, err := NewMemberEntityFromSignUp(signUp)
			if err != nil {
				return err
			}

			return repository.Create(ctx, &newMember)
		}

		return err
	}

	return domain.ErrDuplicated
}

func (MemberService) ApproveMember(ctx context.Context, memberId uint) error {
	repository := memberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	if err := memberEntity.Approve(ctx); err != nil {
		return err
	}

	return repository.Save(ctx, &memberEntity)
}

func (MemberService) GetMemberByGoogleId(ctx context.Context, googleId string) (MemberEntity, error) {
	return memberRepository{}.FindByGoogleId(ctx, googleId)
}

func (MemberService) RejectMember(ctx context.Context, memberId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := memberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdatedBy = userClaim.Id
	return repository.Delete(ctx, memberEntity)
}

func (MemberService) UpdateMemberLastAccessAt(ctx context.Context, memberId uint) error {
	repository := memberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdateLastAccessAt()

	return repository.Save(ctx, &memberEntity)
}
