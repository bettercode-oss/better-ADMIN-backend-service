package services

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	memberDomain "better-admin-backend-service/member/domain"
	"better-admin-backend-service/security"
	"context"
	"github.com/mitchellh/mapstructure"
	pkgerrors "github.com/pkg/errors"
)

type AuthService struct {
	memberService       *MemberService
	organizationService *OrganizationService
	siteService         *SiteService
}

func NewAuthService(
	memberService *MemberService,
	organizationService *OrganizationService,
	siteService *SiteService) *AuthService {

	return &AuthService{
		memberService:       memberService,
		organizationService: organizationService,
		siteService:         siteService,
	}
}

func (s AuthService) AuthWithSignIdPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	memberEntity, err := s.memberService.GetMemberBySignId(ctx, signIn.Id)
	if err != nil {
		return security.JwtToken{}, err
	}

	err = memberEntity.ValidatePassword(signIn.Password)
	if err != nil {
		return security.JwtToken{}, errors.ErrAuthentication
	}

	approved := memberEntity.IsApproved()
	if approved == false {
		return security.JwtToken{}, errors.ErrUnApproved
	}

	return s.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}

func (s AuthService) generateJwtTokenAndLogMemberAccess(ctx context.Context, memberEntity memberDomain.MemberEntity) (token security.JwtToken, err error) {
	memberAssignedAllRoleAndPermission, err := s.organizationService.GetMemberAssignedAllRoleAndPermission(ctx, memberEntity)
	if err != nil {
		return
	}

	token, err = security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
		Id:          memberEntity.ID,
		Roles:       memberAssignedAllRoleAndPermission.Roles,
		Permissions: memberAssignedAllRoleAndPermission.Permissions,
	})

	err = s.logMemberAccessAt(ctx, memberEntity.ID)
	return
}

func (s AuthService) logMemberAccessAt(ctx context.Context, memberId uint) error {
	err := s.memberService.UpdateMemberLastAccessAt(ctx, memberId)
	if err != nil {
		return err
	}

	return nil
}

func (s AuthService) AuthWithDoorayIdAndPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	doorayLoginSetting, err := s.siteService.GetSettingWithKey(ctx, constants.SettingKeyDoorayLogin)
	if err != nil {
		return security.JwtToken{}, err
	}

	var settings dtos.DoorayLoginSetting
	if err = mapstructure.Decode(doorayLoginSetting, &settings); err != nil {
		return security.JwtToken{}, err
	}

	if *settings.Used == false {
		err = pkgerrors.New("not supported dooray login")
		return security.JwtToken{}, err
	}

	doorayMember, err := adapters.DoorayAdapter{}.Authenticate(settings.Domain, settings.AuthorizationToken, signIn.Id, signIn.Password)
	if err != nil {
		return security.JwtToken{}, err
	}

	memberEntity, err := s.memberService.GetMemberByDoorayId(ctx, doorayMember.Id)
	if err != nil {
		if err == errors.ErrNotFound {
			newMemberEntity := memberDomain.NewMemberEntityFromDoorayMember(doorayMember)

			if err = s.memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := s.organizationService.GetMemberAssignedAllRoleAndPermission(ctx, newMemberEntity)
			if err != nil {
				return security.JwtToken{}, err
			}

			return security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
				Id:          newMemberEntity.ID,
				Roles:       memberAssignedAllRoleAndPermission.Roles,
				Permissions: memberAssignedAllRoleAndPermission.Permissions,
			})
		}
		return security.JwtToken{}, err
	}

	return s.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}

func (s AuthService) AuthWithGoogleWorkspaceAccount(ctx context.Context, code string) (security.JwtToken, error) {
	googleWorkspaceLoginSetting, err := s.siteService.GetSettingWithKey(ctx, constants.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		return security.JwtToken{}, err
	}

	var settings dtos.GoogleWorkspaceLoginSetting
	if err = mapstructure.Decode(googleWorkspaceLoginSetting, &settings); err != nil {
		return security.JwtToken{}, err
	}

	if *settings.Used == false {
		err = pkgerrors.New("not supported google workspace login")
		return security.JwtToken{}, err
	}

	googleMember, err := adapters.GoogleOAuthAdapter{}.Authenticate(code, settings)

	if err != nil {
		return security.JwtToken{}, err
	}

	if googleMember.Hd != settings.Domain {
		return security.JwtToken{}, &errors.ErrInvalidGoogleWorkspaceAccount{
			Domain: settings.Domain,
		}
	}

	memberEntity, err := s.memberService.GetMemberByGoogleId(ctx, googleMember.Id)
	if err != nil {
		if err == errors.ErrNotFound {
			newMemberEntity := memberDomain.NewMemberEntityFromGoogleMember(googleMember)

			if err = s.memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := s.organizationService.GetMemberAssignedAllRoleAndPermission(ctx, newMemberEntity)
			if err != nil {
				return security.JwtToken{}, err
			}

			return security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
				Id:          newMemberEntity.ID,
				Roles:       memberAssignedAllRoleAndPermission.Roles,
				Permissions: memberAssignedAllRoleAndPermission.Permissions,
			})
		}
		return security.JwtToken{}, err
	}

	return s.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}
