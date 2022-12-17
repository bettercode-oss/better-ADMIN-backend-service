package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/middlewares"
	"better-admin-backend-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WebHookController struct {
}

func (controller WebHookController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/web-hooks")
	route.POST("", middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		controller.createWebHook)
	route.GET("", middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		controller.getWebHooks)
	route.GET("/:id", middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		middlewares.HttpEtagCache(0),
		controller.getWebHook)
	route.DELETE("/:id", middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		controller.deleteWebHook)
	route.PUT("/:id", middlewares.PermissionChecker([]string{domain.PermissionManageSystemSettings}),
		controller.updateWebHook)
	route.POST("/:id/note", middlewares.PermissionChecker([]string{domain.PermissionNoteWebHooks}),
		controller.noteMessage)
}

func (WebHookController) createWebHook(ctx *gin.Context) {
	var webHookInformation dtos.WebHookInformation
	if err := ctx.Bind(&webHookInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := services.WebHookService{}.CreateWebHook(ctx.Request.Context(), webHookInformation)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (WebHookController) getWebHooks(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	entities, totalCount, err := services.WebHookService{}.GetWebHooks(ctx.Request.Context(), pageable)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
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

	ctx.JSON(http.StatusOK, pageResult)
}

func (WebHookController) getWebHook(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entity, err := services.WebHookService{}.GetWebHook(ctx.Request.Context(), uint(webHookId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	webHookDetails := dtos.WebHookDetails{
		Id:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
	}

	webHookDetails.FillInWebHookCallSpec(ctx.Request, entity.AccessToken)

	ctx.JSON(http.StatusOK, webHookDetails)
}

func (WebHookController) deleteWebHook(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.WebHookService{}.DeleteWebHook(ctx.Request.Context(), uint(webHookId))
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (WebHookController) updateWebHook(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var webHookInformation dtos.WebHookInformation
	if err := ctx.Bind(&webHookInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.WebHookService{}.UpdateWebHook(ctx.Request.Context(), uint(webHookId), webHookInformation)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (WebHookController) noteMessage(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var message dtos.WebHookMessage
	if err := ctx.BindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.WebHookService{}.NoteMessage(ctx.Request.Context(), uint(webHookId), message)
	if err != nil {
		if err == domain.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}
