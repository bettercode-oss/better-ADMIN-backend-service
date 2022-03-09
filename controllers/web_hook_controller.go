package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/webhook"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type WebHookController struct {
}

func (controller WebHookController) Init(g *echo.Group) {
	g.POST("", controller.CreateWebHook, middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.GET("", controller.GetWebHooks, middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.GET("/:id", controller.GetWebHook, middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.DELETE("/:id", controller.DeleteWebHook, middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.PUT("/:id", controller.UpdateWebHook, middlewares.CheckPermission([]string{domain.PermissionManageSystemSettings}))
	g.POST("/:id/note", controller.NoteMessage, middlewares.CheckPermission([]string{domain.PermissionNoteWebHooks}))
}

func (WebHookController) CreateWebHook(ctx echo.Context) error {
	var webHookInformation dtos.WebHookInformation
	if err := ctx.Bind(&webHookInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := webHookInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := webhook.WebHookService{}.CreateWebHook(ctx.Request().Context(), webHookInformation)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (WebHookController) GetWebHooks(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)

	entities, totalCount, err := webhook.WebHookService{}.GetWebHooks(ctx.Request().Context(), pageable)
	if err != nil {
		return err
	}

	var webHooks = make([]dtos.WebHookInformation, 0)
	for _, entity := range entities {
		webHooks = append(webHooks, dtos.WebHookInformation{
			Id:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
		})
	}

	pageResult := dtos.PageResult{
		Result:     webHooks,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, pageResult)
}

func (WebHookController) DeleteWebHook(ctx echo.Context) error {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = webhook.WebHookService{}.DeleteWebHook(ctx.Request().Context(), uint(webHookId))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (WebHookController) GetWebHook(ctx echo.Context) error {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	entity, err := webhook.WebHookService{}.GetWebHook(ctx.Request().Context(), uint(webHookId))
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		return err
	}

	webHookDetails := dtos.WebHookDetails{
		Id:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
	}

	webHookDetails.FillInWebHookCallSpec(ctx.Request(), entity.AccessToken)

	return ctx.JSON(http.StatusOK, webHookDetails)
}

func (WebHookController) UpdateWebHook(ctx echo.Context) error {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var webHookInformation dtos.WebHookInformation
	if err := ctx.Bind(&webHookInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := webHookInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = webhook.WebHookService{}.UpdateWebHook(ctx.Request().Context(), uint(webHookId), webHookInformation)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (WebHookController) NoteMessage(ctx echo.Context) error {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var message dtos.WebHookMessage
	if err := ctx.Bind(&message); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := message.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = webhook.WebHookService{}.NoteMessage(ctx.Request().Context(), uint(webHookId), message)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
