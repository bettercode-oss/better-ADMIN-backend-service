package rest

import (
	memberRepository "better-admin-backend-service/member/repository"
	organizationRepository "better-admin-backend-service/organization/repository"
	rbacRepository "better-admin-backend-service/rbac/repository"
	"better-admin-backend-service/services"
	siteRepository "better-admin-backend-service/site/repository"
	webHookRepository "better-admin-backend-service/webhook/repository"
	"github.com/gin-gonic/gin"
)

type Router struct {
}

func (Router) MapRoutes(routerGroup *gin.RouterGroup) {
	rbacService := services.NewRoleBasedAccessControlService(&rbacRepository.PermissionRepository{}, &rbacRepository.RoleRepository{})
	memberService := services.NewMemberService(rbacService, &memberRepository.MemberRepository{})
	organizationService := services.NewOrganizationService(rbacService, &organizationRepository.OrganizationRepository{}, memberService)
	siteService := services.NewSiteService(&siteRepository.SiteSettingRepository{})
	webHookService := services.NewWebHookService(&webHookRepository.WebHookRepository{})
	authService := services.NewAuthService(memberService, organizationService, siteService)

	NewAccessControlController(
		routerGroup,
		rbacService,
	).MapRoutes()

	NewMemberController(
		routerGroup,
		rbacService,
		memberService,
		organizationService,
	).MapRoutes()

	NewOrganizationController(
		routerGroup,
		organizationService,
	).MapRoutes()

	NewSiteController(
		routerGroup,
		siteService,
	).MapRoutes()

	NewWebHookController(
		routerGroup,
		webHookService,
	).MapRoutes()

	NewAuthController(
		routerGroup,
		authService,
		memberService,
	).MapRoutes()
}
