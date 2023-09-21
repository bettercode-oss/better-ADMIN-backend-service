package middlewares

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	jwtAuthentication security.JwtAuthentication
)

func JwtToken() gin.HandlerFunc {
	jwtAuthentication = security.JwtAuthentication{}

	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		if len(accessToken) == 0 {
			c.Next()
			return
		}

		index := strings.Index(accessToken, "Bearer")
		if index < 0 {
			index = strings.Index(accessToken, "Bearer")
		}
		if index >= 0 {
			accessToken = accessToken[index+len("Bearer"):]
			accessToken = strings.Trim(accessToken, " ")
		}

		userClaim, err := jwtAuthentication.ConvertTokenUserClaim(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dtos.ErrorMessage{Message: err.Error()})
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(helpers.ContextHelper().SetUserClaim(c.Request.Context(), userClaim))
		c.Next()
	}
}
