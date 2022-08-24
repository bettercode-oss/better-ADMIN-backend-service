package services

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/member/entity"
	siteEntity "better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/security"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type AuthService struct {
}

func (service AuthService) AuthWithSignIdPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	memberEntity, err := MemberService{}.GetMemberBySignId(ctx, signIn.Id)
	if err != nil {
		return security.JwtToken{}, err
	}

	err = memberEntity.ValidatePassword(signIn.Password)
	if err != nil {
		return security.JwtToken{}, domain.ErrAuthentication
	}

	approved := memberEntity.IsApproved()
	if approved == false {
		return security.JwtToken{}, domain.ErrUnApproved
	}

	return service.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}

func (service AuthService) generateJwtTokenAndLogMemberAccess(ctx context.Context, memberEntity entity.MemberEntity) (token security.JwtToken, err error) {
	memberAssignedAllRoleAndPermission, err := OrganizationService{}.GetMemberAssignedAllRoleAndPermission(ctx, memberEntity)
	if err != nil {
		return
	}

	token, err = security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
		Id:          memberEntity.ID,
		Roles:       memberAssignedAllRoleAndPermission.Roles,
		Permissions: memberAssignedAllRoleAndPermission.Permissions,
	})

	err = service.logMemberAccessAt(ctx, memberEntity.ID)
	return
}

func (service AuthService) logMemberAccessAt(ctx context.Context, memberId uint) error {
	err := MemberService{}.UpdateMemberLastAccessAt(ctx, memberId)
	if err != nil {
		return err
	}

	return nil
}

func (service AuthService) AuthWithDoorayIdAndPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	doorayLoginSetting, err := SiteService{}.GetSettingWithKey(ctx, siteEntity.SettingKeyDoorayLogin)
	if err != nil {
		return security.JwtToken{}, err
	}

	var settings dtos.DoorayLoginSetting
	if err = mapstructure.Decode(doorayLoginSetting, &settings); err != nil {
		return security.JwtToken{}, err
	}

	if *settings.Used == false {
		err = errors.New("not supported dooray login")
		return security.JwtToken{}, err
	}

	doorayMember, err := adapters.DoorayAdapter{}.Authenticate(settings.Domain, settings.AuthorizationToken, signIn.Id, signIn.Password)
	if err != nil {
		return security.JwtToken{}, err
	}

	memberService := MemberService{}
	memberEntity, err := memberService.GetMemberByDoorayId(ctx, doorayMember.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			newMemberEntity := entity.NewMemberEntityFromDoorayMember(doorayMember)

			if err = memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := OrganizationService{}.GetMemberAssignedAllRoleAndPermission(ctx, newMemberEntity)
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

	return service.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}

func (service AuthService) AuthWithGoogleWorkspaceAccount(ctx context.Context, code string) (security.JwtToken, error) {
	googleWorkspaceLoginSetting, err := SiteService{}.GetSettingWithKey(ctx, siteEntity.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		return security.JwtToken{}, err
	}

	var settings dtos.GoogleWorkspaceLoginSetting
	if err = mapstructure.Decode(googleWorkspaceLoginSetting, &settings); err != nil {
		return security.JwtToken{}, err
	}

	if *settings.Used == false {
		err = errors.New("not supported google workspace login")
		return security.JwtToken{}, err
	}

	googleMember, err := adapters.GoogleOAuthAdapter{}.Authenticate(code, settings)

	if err != nil {
		return security.JwtToken{}, err
	}

	if googleMember.Hd != settings.Domain {
		return security.JwtToken{}, &domain.ErrInvalidGoogleWorkspaceAccount{
			Domain: settings.Domain,
		}
	}

	memberService := MemberService{}
	memberEntity, err := memberService.GetMemberByGoogleId(ctx, googleMember.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			newMemberEntity := entity.NewMemberEntityFromGoogleMember(googleMember)

			if err = memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := OrganizationService{}.GetMemberAssignedAllRoleAndPermission(ctx, newMemberEntity)
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

	return service.generateJwtTokenAndLogMemberAccess(ctx, memberEntity)
}
