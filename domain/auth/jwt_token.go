package auth

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/domain/member"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtTokenGenerator struct {
	member member.MemberEntity
}

func (g JwtTokenGenerator) Generate() (map[string]string, error) {
	accessTokeClaims := jwt.MapClaims{
		"id":          g.member.ID,
		"roles":       "",
		"permissions": "",
		"iss":         "better-admin",
		"aud":         "better-admin",
		"exp":         time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokeClaims).SignedString([]byte(config.Config.JwtSecret))

	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{
		"id":  g.member.ID,
		"iss": "better-admin",
		"aud": "better-admin",
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(config.Config.JwtSecret))

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}, nil
}
