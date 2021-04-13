package auth

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/security"
	"context"
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

	token, err = security.JwtAuthentication{}.GenerateJwtToken(security.UserClaim{
		Id: memberEntity.ID,
	})
	return
}
