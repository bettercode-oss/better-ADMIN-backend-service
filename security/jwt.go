package security

import (
	"better-admin-backend-service/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

// https://docs.apigee.com/api-platform/reference/policies/oauth-http-status-code-reference
var InvalidAccessToken = errors.New("invalid access token")
var AccessTokenExpired = errors.New("access token expired")

type JwtAuthentication struct {
}

func (JwtAuthentication) GenerateJwtToken(claim UserClaim) (JwtToken, error) {
	claimMap, err := claim.ConvertMap()
	if err != nil {
		return JwtToken{}, err
	}

	accessTokenClaims := jwt.MapClaims{}
	for key, value := range claimMap {
		accessTokenClaims[key] = value
	}
	// TODO config 로 빼자.
	accessTokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(config.Config.JwtSecret))

	if err != nil {
		return JwtToken{}, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	for key, value := range claimMap {
		refreshTokenClaims[key] = value
	}
	// TODO config 로 빼자.
	refreshTokenClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(config.Config.JwtSecret))

	return JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (JwtAuthentication) GenerateJwtAccessTokenNeverExpired(claim UserClaim) (string, error) {
	claimMap, err := claim.ConvertMap()
	if err != nil {
		return "", err
	}

	accessTokenClaims := jwt.MapClaims{}
	for key, value := range claimMap {
		accessTokenClaims[key] = value
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(config.Config.JwtSecret))

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (JwtAuthentication) ConvertTokenUserClaim(token string) (*UserClaim, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(config.Config.JwtSecret), nil })

	if err != nil {
		log.Error("JWT parsing error: " + err.Error())
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, AccessTokenExpired
			}
		}
		return nil, InvalidAccessToken
	}

	if jwt.SigningMethodHS256.Alg() != parsedToken.Header["alg"] {
		log.Error(fmt.Sprintf("Error: jwt token is expected %s signing method but token specified %s",
			jwt.SigningMethodHS256.Alg(), parsedToken.Header["alg"]))
		return nil, InvalidAccessToken
	}

	if !parsedToken.Valid {
		return nil, InvalidAccessToken
	}

	claimInfo, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Can'get jwt.MapClaims")
		return nil, InvalidAccessToken
	}

	userClaim, err := NewUserClaim(claimInfo)
	if err != nil {
		return nil, err
	}

	return &userClaim, nil
}

func (jwtAuthentication JwtAuthentication) RefreshAccessToken(refreshToken string) (string, error) {
	userClaim, err := jwtAuthentication.ConvertTokenUserClaim(refreshToken)
	if err != nil {
		return "", err
	}

	jwtToken, err := jwtAuthentication.GenerateJwtToken(*userClaim)
	if err != nil {
		return "", err
	}

	return jwtToken.AccessToken, nil
}

func (jwtAuthentication JwtAuthentication) ValidateToken(token string) error {
	_, err := jwtAuthentication.ConvertTokenUserClaim(token)
	return err
}

type JwtToken struct {
	AccessToken  string
	RefreshToken string
}

type UserClaim struct {
	Id          uint     `json:"id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (c UserClaim) ConvertMap() (map[string]interface{}, error) {
	bytes, err := json.Marshal(c)

	if err != nil {
		return nil, err
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(bytes, &resultMap); err != nil {
		return nil, err
	}

	return resultMap, nil
}

func NewUserClaim(mapUserClaim map[string]interface{}) (UserClaim, error) {
	bytes, err := json.Marshal(mapUserClaim)
	if err != nil {
		return UserClaim{}, err
	}

	var claim UserClaim
	if err := json.Unmarshal(bytes, &claim); err != nil {
		return UserClaim{}, err
	}

	return claim, nil
}
