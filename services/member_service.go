package services

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/member/entity"
	"better-admin-backend-service/domain/member/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
)

type MemberService struct {
}

func (MemberService) GetMemberBySignId(ctx context.Context, signId string) (entity.MemberEntity, error) {
	return repository.MemberRepository{}.FindBySignId(ctx, signId)
}

func (MemberService) GetMemberByDoorayId(ctx context.Context, doorayId string) (entity.MemberEntity, error) {
	return repository.MemberRepository{}.FindByDoorayId(ctx, doorayId)
}

func (MemberService) CreateMember(ctx context.Context, entity *entity.MemberEntity) error {
	return repository.MemberRepository{}.Create(ctx, entity)
}

func (MemberService) GetMemberById(ctx context.Context, id uint) (entity.MemberEntity, error) {
	return repository.MemberRepository{}.FindById(ctx, id)
}

func (MemberService) GetMembers(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]entity.MemberEntity, int64, error) {
	return repository.MemberRepository{}.FindAll(ctx, filters, pageable)
}

func (MemberService) AssignRole(ctx context.Context, memberId uint, assignRole dtos.MemberAssignRole) error {
	repository := repository.MemberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["roleIds"] = assignRole.RoleIds

	findRoleEntities, _, err := RoleBasedAccessControlService{}.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = memberEntity.AssignRole(ctx, findRoleEntities)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &memberEntity)
}

func (MemberService) GetMember(ctx context.Context, memberId uint) (entity.MemberEntity, error) {
	return repository.MemberRepository{}.FindById(ctx, memberId)
}

func (MemberService) SignUpMember(ctx context.Context, signUp dtos.MemberSignUp) error {
	repository := repository.MemberRepository{}
	_, err := repository.FindBySignId(ctx, signUp.SignId)
	if err != nil {
		if err == domain.ErrNotFound {
			// signId 가 중복이 없을 때만 가입
			newMember, err := entity.NewMemberEntityFromSignUp(signUp)
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
	repository := repository.MemberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	if err := memberEntity.Approve(ctx); err != nil {
		return err
	}

	return repository.Save(ctx, &memberEntity)
}

func (MemberService) GetMemberByGoogleId(ctx context.Context, googleId string) (entity.MemberEntity, error) {
	return repository.MemberRepository{}.FindByGoogleId(ctx, googleId)
}

func (MemberService) RejectMember(ctx context.Context, memberId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := repository.MemberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdatedBy = userClaim.Id
	return repository.Delete(ctx, memberEntity)
}

func (MemberService) UpdateMemberLastAccessAt(ctx context.Context, memberId uint) error {
	repository := repository.MemberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdateLastAccessAt()

	return repository.Save(ctx, &memberEntity)
}
