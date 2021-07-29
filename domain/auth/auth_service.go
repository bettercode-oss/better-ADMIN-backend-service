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

func (s AuthService) AuthWithSignIdPassword(ctx context.Context, signIn dtos.MemberSignIn) (token security.JwtToken, err error) {
	memberEntity, err := member.MemberService{}.GetMemberBySignId(ctx, signIn.Id)
	if err != nil {
		return
	}

	err = memberEntity.ValidatePassword(signIn.Password)
	if err != nil {
		err = domain.ErrAuthentication
		return
	}

	approved := memberEntity.IsApproved()
	if approved == false {
		err = domain.ErrUnApproved
		return
	}

	memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx, memberEntity)
	if err != nil {
		return
	}

	token, err = security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
		Id:          memberEntity.ID,
		Roles:       memberAssignedAllRoleAndPermission.Roles,
		Permissions: memberAssignedAllRoleAndPermission.Permissions,
	})
	return
}

func (AuthService) AuthWithDoorayIdAndPassword(ctx context.Context, signIn dtos.MemberSignIn) (security.JwtToken, error) {
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

	memberAssignedAllRoleAndPermission, err := factory.MemberAssignedAllRoleAndPermissionFactory{}.Create(ctx, memberEntity)
	if err != nil {
		return security.JwtToken{}, err
	}

	return security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
		Id:          memberEntity.ID,
		Roles:       memberAssignedAllRoleAndPermission.Roles,
		Permissions: memberAssignedAllRoleAndPermission.Permissions,
	})
}
