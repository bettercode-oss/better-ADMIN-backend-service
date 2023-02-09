package rest

import (
	"better-admin-backend-service/app/middlewares"
	"better-admin-backend-service/constants"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/services"
	etag "github.com/bettercode-oss/gin-middleware-etag"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type WebHookController struct {
	routerGroup    *gin.RouterGroup
	webHookService *services.WebHookService
}

func NewWebHookController(
	routerGroup *gin.RouterGroup,
	webHookService *services.WebHookService) *WebHookController {

	return &WebHookController{
		routerGroup:    routerGroup,
		webHookService: webHookService,
	}
}

func (c WebHookController) MapRoutes() {
	route := c.routerGroup.Group("/web-hooks")
	route.POST("", middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		c.createWebHook)
	route.GET("", middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		etag.HttpEtagCache(0),
		c.getWebHooks)
	route.GET("/:id", middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		etag.HttpEtagCache(0),
		c.getWebHook)
	route.DELETE("/:id", middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		c.deleteWebHook)
	route.PUT("/:id", middlewares.PermissionChecker([]string{constants.PermissionManageSystemSettings}),
		c.updateWebHook)
	route.POST("/:id/note", middlewares.PermissionChecker([]string{constants.PermissionNoteWebHooks}),
		c.noteMessage)
}

func (c WebHookController) createWebHook(ctx *gin.Context) {
	var webHookInformation dtos.WebHookInformation
	if err := ctx.Bind(&webHookInformation); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.webHookService.CreateWebHook(ctx.Request.Context(), webHookInformation)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c WebHookController) getWebHooks(ctx *gin.Context) {
	pageable := dtos.NewPageableFromRequest(ctx)

	entities, totalCount, err := c.webHookService.GetWebHooks(ctx.Request.Context(), pageable)
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

func (c WebHookController) getWebHook(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	entity, err := c.webHookService.GetWebHook(ctx.Request.Context(), uint(webHookId))
	if err != nil {
		if err == errors.ErrNotFound {
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

func (c WebHookController) deleteWebHook(ctx *gin.Context) {
	webHookId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.webHookService.DeleteWebHook(ctx.Request.Context(), uint(webHookId))
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c WebHookController) updateWebHook(ctx *gin.Context) {
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

	err = c.webHookService.UpdateWebHook(ctx.Request.Context(), uint(webHookId), webHookInformation)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c WebHookController) noteMessage(ctx *gin.Context) {
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

	err = c.webHookService.NoteMessage(ctx.Request.Context(), uint(webHookId), message)
	if err != nil {
		if err == errors.ErrNotFound {
			ctx.Status(http.StatusNotFound)
			return
		}
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}
