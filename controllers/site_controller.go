package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

type SiteController struct {
}

func (controller SiteController) Init(g *echo.Group) {
	g.GET("/settings", controller.GetSettingsSummary)
	g.PUT("/settings/dooray-login", controller.SetDoorayLoginSetting, middlewares.CheckPermission([]string{"MANAGE_SYSTEM_SETTINGS"}))
	g.GET("/settings/dooray-login", controller.GetDoorayLoginSetting, middlewares.CheckPermission([]string{"MANAGE_SYSTEM_SETTINGS"}))
	g.GET("/settings/google-workspace-login", controller.GetGoogleWorkspaceLoginSetting, middlewares.CheckPermission([]string{"MANAGE_SYSTEM_SETTINGS"}))
	g.PUT("/settings/google-workspace-login", controller.SetGoogleWorkspaceLoginSetting, middlewares.CheckPermission([]string{"MANAGE_SYSTEM_SETTINGS"}))
}

func (controller SiteController) SetDoorayLoginSetting(ctx echo.Context) error {
	var setting dtos.DoorayLoginSetting

	if err := ctx.Bind(&setting); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := setting.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	service := site.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request().Context(), site.SettingKeyDoorayLogin, setting); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller SiteController) GetDoorayLoginSetting(ctx echo.Context) error {
	setting, err := site.SiteService{}.GetSettingWithKey(ctx.Request().Context(), site.SettingKeyDoorayLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusOK, dtos.DoorayLoginSetting{})
		}

		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, setting)
}

func (controller SiteController) GetSettingsSummary(ctx echo.Context) error {
	settings, err := site.SiteService{}.GetSettings(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	summary := dtos.SiteSettingsSummary{}

	for _, setting := range settings {
		if setting.Key == site.SettingKeyDoorayLogin {
			var doorayLoginSetting dtos.DoorayLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &doorayLoginSetting)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			if *doorayLoginSetting.Used {
				summary.DoorayLoginUsed = true
			}
		}

		if setting.Key == site.SettingKeyGoogleWorkspaceLogin {
			var googleWorkspaceSetting dtos.GoogleWorkspaceLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &googleWorkspaceSetting)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			if *googleWorkspaceSetting.Used {
				summary.GoogleWorkspaceLoginUsed = true
				summary.GoogleWorkspaceOAuthUri = googleWorkspaceSetting.GetOAuthUri()
			}
		}
	}
	return ctx.JSON(http.StatusOK, summary)
}

func (SiteController) GetGoogleWorkspaceLoginSetting(ctx echo.Context) error {
	setting, err := site.SiteService{}.GetSettingWithKey(ctx.Request().Context(), site.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusOK, dtos.GoogleWorkspaceLoginSetting{})
		}

		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, setting)
}

func (SiteController) SetGoogleWorkspaceLoginSetting(ctx echo.Context) error {
	var setting dtos.GoogleWorkspaceLoginSetting

	if err := ctx.Bind(&setting); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := setting.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	service := site.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request().Context(), site.SettingKeyGoogleWorkspaceLogin, setting); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}
