package auth

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/factory"
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/security"
	"context"
	"errors"
	"github.com/mitchellh/mapstructure"
)

type AuthService struct {
}

func (service AuthService) AuthWithSignIdPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	memberEntity, err := member.MemberService{}.GetMemberBySignId(ctx, signIn.Id)
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

func (service AuthService) generateJwtTokenAndLogMemberAccess(ctx context.Context, memberEntity member.MemberEntity) (token security.JwtToken, err error) {
	memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx, memberEntity)
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
	err := member.MemberService{}.UpdateMemberLastAccessAt(ctx, memberId)
	if err != nil {
		return err
	}

	return nil
}

func (service AuthService) AuthWithDoorayIdAndPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
	doorayLoginSetting, err := site.SiteService{}.GetSettingWithKey(ctx, site.SettingKeyDoorayLogin)
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

	memberService := member.MemberService{}
	memberEntity, err := memberService.GetMemberByDoorayId(ctx, doorayMember.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			newMemberEntity := member.NewMemberEntityFromDoorayMember(doorayMember)

			if err = memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx, newMemberEntity)
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
	googleWorkspaceLoginSetting, err := site.SiteService{}.GetSettingWithKey(ctx, site.SettingKeyGoogleWorkspaceLogin)
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

	memberService := member.MemberService{}
	memberEntity, err := memberService.GetMemberByGoogleId(ctx, googleMember.Id)
	if err != nil {
		if err == domain.ErrNotFound {
			newMemberEntity := member.NewMemberEntityFromGoogleMember(googleMember)

			if err = memberService.CreateMember(ctx, &newMemberEntity); err != nil {
				return security.JwtToken{}, err
			}

			memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx, newMemberEntity)
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
