package services

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/member/domain"
	"better-admin-backend-service/member/repository"
	"context"
)

type MemberService struct {
	rbacService      *RoleBasedAccessControlService
	memberRepository *repository.MemberRepository
}

func NewMemberService(rbacService *RoleBasedAccessControlService,
	memberRepository *repository.MemberRepository) *MemberService {
	return &MemberService{
		rbacService:      rbacService,
		memberRepository: memberRepository,
	}
}

func (s MemberService) GetMemberBySignId(ctx context.Context, signId string) (domain.MemberEntity, error) {
	return s.memberRepository.FindBySignId(ctx, signId)
}

func (s MemberService) GetMemberByDoorayId(ctx context.Context, doorayId string) (domain.MemberEntity, error) {
	return s.memberRepository.FindByDoorayId(ctx, doorayId)
}

func (s MemberService) CreateMember(ctx context.Context, entity *domain.MemberEntity) error {
	return s.memberRepository.Create(ctx, entity)
}

func (s MemberService) GetMemberById(ctx context.Context, id uint) (domain.MemberEntity, error) {
	return s.memberRepository.FindById(ctx, id)
}

func (s MemberService) GetMembers(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]domain.MemberEntity, int64, error) {
	return s.memberRepository.FindAll(ctx, filters, pageable)
}

func (s MemberService) AssignRole(ctx context.Context, memberId uint, assignRole dtos.MemberAssignRole) error {
	memberEntity, err := s.memberRepository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["roleIds"] = assignRole.RoleIds

	findRoleEntities, _, err := s.rbacService.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = memberEntity.AssignRole(ctx, findRoleEntities)
	if err != nil {
		return err
	}

	return s.memberRepository.Save(ctx, &memberEntity)
}

func (s MemberService) GetMember(ctx context.Context, memberId uint) (domain.MemberEntity, error) {
	return s.memberRepository.FindById(ctx, memberId)
}

func (s MemberService) SignUpMember(ctx context.Context, signUp dtos.MemberSignUp) error {
	_, err := s.memberRepository.FindBySignId(ctx, signUp.SignId)
	if err != nil {
		if err == errors.ErrNotFound {
			// signId 가 중복이 없을 때만 가입
			newMember, err := domain.NewMemberEntityFromSignUp(signUp)
			if err != nil {
				return err
			}

			return s.memberRepository.Create(ctx, &newMember)
		}

		return err
	}

	return errors.ErrDuplicated
}

func (s MemberService) ApproveMember(ctx context.Context, memberId uint) error {
	memberEntity, err := s.memberRepository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	if err := memberEntity.Approve(ctx); err != nil {
		return err
	}

	return s.memberRepository.Save(ctx, &memberEntity)
}

func (s MemberService) GetMemberByGoogleId(ctx context.Context, googleId string) (domain.MemberEntity, error) {
	return s.memberRepository.FindByGoogleId(ctx, googleId)
}

func (s MemberService) RejectMember(ctx context.Context, memberId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	memberEntity, err := s.memberRepository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdatedBy = userClaim.Id
	return s.memberRepository.Delete(ctx, memberEntity)
}

func (s MemberService) UpdateMemberLastAccessAt(ctx context.Context, memberId uint) error {
	memberEntity, err := s.memberRepository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	memberEntity.UpdateLastAccessAt()

	return s.memberRepository.Save(ctx, &memberEntity)
}
