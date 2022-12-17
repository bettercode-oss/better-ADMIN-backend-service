package middlewares

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func PermissionChecker(allowPermissions []string) gin.HandlerFunc {
	allowPermissionMap := make(map[string]bool)
	for _, permission := range allowPermissions {
		allowPermissionMap[permission] = true
	}

	return func(ctx *gin.Context) {
		userClaim, err := helpers.ContextHelper().GetUserClaim(ctx.Request.Context())
		if err != nil {
			log.Warnf("No valid credentials: %s", ctx.Request.RequestURI)
			ctx.JSON(http.StatusUnauthorized, "Please provide valid credentials")
			ctx.Abort()
			return
		}
		if len(allowPermissions) == 1 && allowPermissions[0] == "*" {
			ctx.Next()
			return
		} else {
			for _, permission := range userClaim.Permissions {
				if allowPermissionMap[permission] {
					ctx.Next()
					return
				}
			}
			log.Warnf("Can't access this API: %s", ctx.Request.RequestURI)
			ctx.JSON(http.StatusForbidden, "Can't access this API")
			ctx.Abort()
			return
		}
	}
}
