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

func (g JwtTokenGenerator) Generate() (string, error) {
	claims := jwt.MapClaims{
		"id":          g.member.ID,
		"roles":       "",
		"permissions": "",
		"iss":         "better-admin",
		"aud":         "better-admin",
		"nbf":         time.Now().Add(-time.Minute * 5).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Config.JwtSecret))
	return token, err
}
