package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(e *gin.Engine) {
	apiRoute := e.Group("/api")
	AuthController{}.Init(apiRoute)
	SiteController{}.Init(apiRoute)
	MemberController{}.Init(apiRoute)
	AccessControlController{}.Init(apiRoute)
	OrganizationController{}.Init(apiRoute)
	WebHookController{}.Init(apiRoute)
}
