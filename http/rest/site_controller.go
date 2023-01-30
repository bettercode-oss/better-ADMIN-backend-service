package rest

import (
	"better-admin-backend-service/app/middlewares"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	pkgerrors "github.com/pkg/errors"
	"net/http"
)

type SiteController struct {
	routerGroup *gin.RouterGroup
	siteService *services.SiteService
}

func NewSiteController(
	routerGroup *gin.RouterGroup,
	siteService *services.SiteService) *SiteController {

	return &SiteController{
		routerGroup: routerGroup,
		siteService: siteService,
	}
}

func (c SiteController) MapRoutes() {
	route := c.routerGroup.Group("/site")

	route.GET("/settings",
		middlewares.HttpEtagCache(0),
		c.getSettingsSummary)
	route.GET("/settings/dooray-login",
		middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		c.getDoorayLoginSetting)
	route.PUT("/settings/dooray-login",
		middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		c.setDoorayLoginSetting)
	route.GET("/settings/google-workspace-login",
		middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		c.getGoogleWorkspaceLoginSetting)
	route.PUT("/settings/google-workspace-login",
		middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		c.setGoogleWorkspaceLoginSetting)
	route.GET("/settings/app-version",
		middlewares.HttpEtagCache(0),
		c.getAppVersion)
	route.PUT("/settings/app-version",
		c.increaseAppVersion)
}
func (c SiteController) getSettingsSummary(ctx *gin.Context) {
	settings, err := c.siteService.GetSettings(ctx.Request.Context())

	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	summary := dtos.SiteSettingsSummary{}

	for _, setting := range settings {
		if setting.Key == constants.SettingKeyDoorayLogin {
			var doorayLoginSetting dtos.DoorayLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &doorayLoginSetting)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, pkgerrors.Wrap(err, "map to struct decode error"))
				return
			}

			if *doorayLoginSetting.Used {
				summary.DoorayLoginUsed = true
			}
		}

		if setting.Key == constants.SettingKeyGoogleWorkspaceLogin {
			var googleWorkspaceSetting dtos.GoogleWorkspaceLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &googleWorkspaceSetting)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, pkgerrors.Wrap(err, "map to struct decode error"))
				return
			}

			if *googleWorkspaceSetting.Used {
				summary.GoogleWorkspaceLoginUsed = true
				summary.GoogleWorkspaceOAuthUri = googleWorkspaceSetting.GetOAuthUri()
			}
		}
	}

	ctx.JSON(http.StatusOK, summary)
}

func (c SiteController) getDoorayLoginSetting(ctx *gin.Context) {
	setting, err := c.siteService.GetSettingWithKey(ctx.Request.Context(), constants.SettingKeyDoorayLogin)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.JSON(http.StatusOK, dtos.DoorayLoginSetting{})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, setting)
}

func (c SiteController) setDoorayLoginSetting(ctx *gin.Context) {
	var setting dtos.DoorayLoginSetting

	if err := ctx.BindJSON(&setting); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := c.siteService.SetSettingWithKey(ctx.Request.Context(), constants.SettingKeyDoorayLogin, setting); err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c SiteController) getGoogleWorkspaceLoginSetting(ctx *gin.Context) {
	setting, err := c.siteService.GetSettingWithKey(ctx.Request.Context(), constants.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.JSON(http.StatusOK, dtos.GoogleWorkspaceLoginSetting{})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, setting)
}

func (c SiteController) setGoogleWorkspaceLoginSetting(ctx *gin.Context) {
	var setting dtos.GoogleWorkspaceLoginSetting

	if err := ctx.BindJSON(&setting); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := c.siteService.SetSettingWithKey(ctx.Request.Context(), constants.SettingKeyGoogleWorkspaceLogin, setting); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c SiteController) getAppVersion(ctx *gin.Context) {
	appVersion, err := c.siteService.GetAppVersion(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, appVersion)
}

func (c SiteController) increaseAppVersion(ctx *gin.Context) {
	err := c.siteService.IncreaseAppVersion(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}
