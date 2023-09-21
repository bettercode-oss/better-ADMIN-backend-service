package middlewares

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/rego"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RestAuthorizer(regoQuery *rego.PreparedEvalQuery) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		input := map[string]any{}
		userClaim, err := helpers.ContextHelper().GetUserClaim(c.Request.Context())
		if err == nil {
			input["member"] = map[string]any{
				"id":          userClaim.Id,
				"permissions": userClaim.Permissions,
			}
		}
		url := c.FullPath()
		if len(url) == 0 {
			url = c.Request.RequestURI
		}

		input["api"] = map[string]any{
			"url":    url,
			"method": c.Request.Method,
		}

		rs, err := regoQuery.Eval(context.TODO(), rego.EvalInput(input))
		if err != nil {
			log.Error("opa error", err)
			c.JSON(http.StatusInternalServerError, dtos.ErrorMessage{Message: err.Error()})
			c.Abort()
			return
		}

		if len(rs) > 0 && len(rs[0].Expressions) > 0 && rs[0].Expressions[0].String() == "true" {
			// 인가
			c.Next()
		} else {
			if userClaim == nil {
				c.JSON(http.StatusUnauthorized, dtos.ErrorMessage{Message: "You are not authorized."})
			} else {
				c.JSON(http.StatusForbidden, dtos.ErrorMessage{Message: "You are not authorized."})
			}
			c.Abort()
			return
		}
	}
}
