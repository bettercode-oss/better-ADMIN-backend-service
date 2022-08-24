package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/site/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/labstack/echo"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net/http"
)

type SiteController struct {
}

func (controller SiteController) Init(g *echo.Group) {
	g.GET("/settings", controller.GetSettingsSummary)
	g.PUT("/settings/dooray-login", controller.SetDoorayLoginSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.GET("/settings/dooray-login", controller.GetDoorayLoginSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.PUT("/settings/google-workspace-login", controller.SetGoogleWorkspaceLoginSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.GET("/settings/google-workspace-login", controller.GetGoogleWorkspaceLoginSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.PUT("/settings/member-access-logs", controller.SetMemberAccessLogSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.GET("/settings/member-access-logs", controller.GetMemberAccessLogSetting,
		middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
}

func (controller SiteController) SetDoorayLoginSetting(ctx echo.Context) error {
	var setting dtos.DoorayLoginSetting

	if err := ctx.Bind(&setting); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := setting.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	service := services.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request().Context(), entity.SettingKeyDoorayLogin, setting); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (controller SiteController) GetDoorayLoginSetting(ctx echo.Context) error {
	setting, err := services.SiteService{}.GetSettingWithKey(ctx.Request().Context(), entity.SettingKeyDoorayLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusOK, dtos.DoorayLoginSetting{})
		}

		return err
	}

	return ctx.JSON(http.StatusOK, setting)
}

func (controller SiteController) GetSettingsSummary(ctx echo.Context) error {
	settings, err := services.SiteService{}.GetSettings(ctx.Request().Context())

	if err != nil {
		return err
	}

	summary := dtos.SiteSettingsSummary{}

	for _, setting := range settings {
		if setting.Key == entity.SettingKeyDoorayLogin {
			var doorayLoginSetting dtos.DoorayLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &doorayLoginSetting)
			if err != nil {
				return errors.Wrap(err, "map to struct decode error")
			}

			if *doorayLoginSetting.Used {
				summary.DoorayLoginUsed = true
			}
		}

		if setting.Key == entity.SettingKeyGoogleWorkspaceLogin {
			var googleWorkspaceSetting dtos.GoogleWorkspaceLoginSetting
			err := mapstructure.Decode(setting.ValueObject, &googleWorkspaceSetting)
			if err != nil {
				return errors.Wrap(err, "map to struct decode error")
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
	setting, err := services.SiteService{}.GetSettingWithKey(ctx.Request().Context(), entity.SettingKeyGoogleWorkspaceLogin)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusOK, dtos.GoogleWorkspaceLoginSetting{})
		}

		return err
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

	service := services.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request().Context(), entity.SettingKeyGoogleWorkspaceLogin, setting); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (SiteController) SetMemberAccessLogSetting(ctx echo.Context) error {
	var setting dtos.MemberAccessLogSetting

	if err := ctx.Bind(&setting); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := setting.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	service := services.SiteService{}
	if err := service.SetSettingWithKey(ctx.Request().Context(), entity.SettingKeyMemberAccessLog, setting); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (SiteController) GetMemberAccessLogSetting(ctx echo.Context) error {
	setting, err := services.SiteService{}.GetSettingWithKey(ctx.Request().Context(), entity.SettingKeyMemberAccessLog)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusOK, dtos.MemberAccessLogSetting{})
		}

		return err
	}

	return ctx.JSON(http.StatusOK, setting)
}
