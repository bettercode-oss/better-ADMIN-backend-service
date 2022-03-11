package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/logging"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"net/http"
)

type MemberAccessLogController struct {
}

func (controller MemberAccessLogController) Init(g *echo.Group) {
	g.POST("", controller.LogMemberAccess, middlewares.CheckPermission([]string{"*"}))
	g.GET("", controller.GetMemberAccessLogs, middlewares.CheckPermission([]string{domain.PermissionViewMonitoring}))
}

func (MemberAccessLogController) LogMemberAccess(ctx echo.Context) error {
	var accessLog dtos.MemberAccessLog

	if err := ctx.Bind(&accessLog); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := accessLog.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	accessLog.IpAddress = ReadUserIP(ctx.Request())
	accessLog.BrowserUserAgent = ctx.Request().UserAgent()

	err := logging.MemberAccessLogService{}.LogMemberAccess(ctx.Request().Context(), accessLog)
	if err != nil {
		if err == domain.ErrNotSupportedAccessLogType {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		return err
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func (MemberAccessLogController) GetMemberAccessLogs(ctx echo.Context) error {
	pageable := dtos.GetPageableFromRequest(ctx)
	filters := map[string]interface{}{}

	if len(ctx.QueryParam("memberId")) > 0 {
		filters["memberId"] = ctx.QueryParam("memberId")
	}

	accessLogEntities, totalCount, err := logging.MemberAccessLogService{}.GetMemberAccessLogs(ctx.Request().Context(), filters, pageable)
	if err != nil {
		return err
	}

	var accessLogs = make([]dtos.MemberAccessLog, 0)
	for _, entity := range accessLogEntities {
		accessLogs = append(accessLogs, dtos.MemberAccessLog{
			Id:               entity.ID,
			MemberId:         entity.MemberId,
			Type:             entity.Type,
			TypeName:         entity.GetType(),
			Url:              entity.Url,
			Method:           entity.Method,
			Parameters:       entity.Parameters,
			Payload:          entity.Payload,
			StatusCode:       entity.StatusCode,
			IpAddress:        entity.IpAddress,
			BrowserUserAgent: entity.BrowserUserAgent,
			CreatedAt:        entity.CreatedAt,
		})
	}

	pageResult := dtos.PageResult{
		Result:     accessLogs,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, pageResult)
}
