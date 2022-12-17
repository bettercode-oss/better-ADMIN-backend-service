package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net/http"
)

type SiteController struct {
}

func (controller SiteController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/site")

	route.GET("/settings",
		middlewares.HttpEtagCache(0),
		controller.getSettingsSummary)
	route.GET("/settings/dooray-login",
		middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		controller.getDoorayLoginSetting)
	route.PUT("/settings/dooray-login",
		middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		controller.setDoorayLoginSetting)
	route.GET("/settings/google-workspace-login",
		middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		controller.getGoogleWorkspaceLoginSetting)
	route.PUT("/settings/google-workspace-login",
		middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		controller.setGoogleWorkspaceLoginSetting)
	route.GET("/settings/app-version",
		middlewares.HttpEtagCache(0),
		controller.getAppVersion)
	route.PUT("/settings/app-version",
		controller.increaseAppVersion)
}
func (SiteController) getSettingsSummary(ctx *gin.Context) {
	settings, err := services.SiteService{}.GetSettings(ctx.Request.Context())

	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	summary := dtos.SiteSettingsSummary{}

	for _, setting := range settings {
		if setting.Key == entity.SettingKeyDoorayLogin {
			var doorayLoginSetting dtos.DoorayLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &doorayLoginSetting)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errors.Wrap(err, "map to struct decode error"))
				return
			}

			if *doorayLoginSetting.Used {
				summary.DoorayLoginUsed = true
			}
		}

		if setting.Key == entity.SettingKeyGoogleWorkspaceLogin {
			var googleWorkspaceSetting dtos.GoogleWorkspaceLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &googleWorkspaceSetting)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errors.Wrap(err, "map to struct decode error"))
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

func (SiteController) getDoorayLoginSetting(ctx *gin.Context) {
	setting, err := services.SiteService{}.GetSettingWithKey(ctx.Request.Context(), entity.SettingKeyDoorayLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.JSON(http.StatusOK, dtos.DoorayLoginSetting{})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, setting)
}

func (SiteController) setDoorayLoginSetting(ctx *gin.Context) {
	var setting dtos.DoorayLoginSetting

	if err := ctx.BindJSON(&setting); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	service := services.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request.Context(), entity.SettingKeyDoorayLogin, setting); err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (SiteController) getGoogleWorkspaceLoginSetting(ctx *gin.Context) {
	setting, err := services.SiteService{}.GetSettingWithKey(ctx.Request.Context(), entity.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.JSON(http.StatusOK, dtos.GoogleWorkspaceLoginSetting{})
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, setting)
}

func (SiteController) setGoogleWorkspaceLoginSetting(ctx *gin.Context) {
	var setting dtos.GoogleWorkspaceLoginSetting

	if err := ctx.BindJSON(&setting); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	service := services.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request.Context(), entity.SettingKeyGoogleWorkspaceLogin, setting); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (SiteController) getAppVersion(ctx *gin.Context) {
	appVersion, err := services.SiteService{}.GetAppVersion(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, appVersion)
}

func (SiteController) increaseAppVersion(ctx *gin.Context) {
	err := services.SiteService{}.IncreaseAppVersion(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}
