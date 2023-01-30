package rest

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/testdata/testserver"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"time"
)

var (
	gormDB *gorm.DB
	ginApp *gin.Engine
)

func init() {
	err := config.InitConfig("../../config/config.json")
	if err != nil {
		panic(err)
	}

	testAppServer := testserver.NewTestAppServer(Router{})
	gormDB = testAppServer.GetDB()
	ginApp = testAppServer.GetGin()
}

func generateTestJWT(claim map[string]interface{}, duration time.Duration) (string, error) {
	token := jwt.MapClaims{}
	for key, value := range claim {
		token[key] = value
	}

	token["exp"] = time.Now().Add(duration).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(config.Config.JwtSecret))
}
