package rest

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"better-admin-backend-service/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	routerGroup   *gin.RouterGroup
	authService   *services.AuthService
	memberService *services.MemberService
}

func NewAuthController(
	routerGroup *gin.RouterGroup,
	authService *services.AuthService,
	memberService *services.MemberService) *AuthController {

	return &AuthController{
		routerGroup:   routerGroup,
		authService:   authService,
		memberService: memberService,
	}
}

func (c AuthController) MapRoutes() {
	route := c.routerGroup.Group("/auth")

	route.POST("", c.authWithSignIdPassword)
	route.POST("/dooray", c.authWithDoorayIdPassword)
	route.GET("/google-workspace", c.authWithGoogleWorkspaceAccount)
	route.POST("/token/refresh", c.refreshAccessToken)
}

func (c AuthController) authWithSignIdPassword(ctx *gin.Context) {
	var memberSignIn dtos.MemberSignIn

	if err := ctx.BindJSON(&memberSignIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := c.authService.AuthWithSignIdPassword(ctx.Request.Context(), memberSignIn)
	if err != nil {
		if err == errors.ErrNotFound || err == errors.ErrAuthentication {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err == errors.ErrUnApproved {
			ctx.JSON(http.StatusNotAcceptable, err.Error())
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	result := map[string]any{}
	result["accessToken"] = jwtToken.AccessToken
	result["expiresAt"] = jwtToken.ExpiresAt
	result["refreshToken"] = jwtToken.RefreshToken
	result["refreshTokenExpiresIn"] = jwtToken.RefreshTokenExpiresIn

	ctx.JSON(http.StatusOK, result)
}

func (c AuthController) authWithDoorayIdPassword(ctx *gin.Context) {
	var memberSignIn dtos.MemberSignIn

	if err := ctx.BindJSON(&memberSignIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := c.authService.AuthWithDoorayIdAndPassword(ctx.Request.Context(), memberSignIn)
	if err != nil {
		if err == errors.ErrAuthentication {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	result := map[string]any{}
	result["accessToken"] = jwtToken.AccessToken
	result["expiresAt"] = jwtToken.ExpiresAt
	result["refreshToken"] = jwtToken.RefreshToken
	result["refreshTokenExpiresIn"] = jwtToken.RefreshTokenExpiresIn

	ctx.JSON(http.StatusOK, result)
}

func (c AuthController) authWithGoogleWorkspaceAccount(ctx *gin.Context) {
	code := ctx.Query("code")
	redirect := ctx.Query("state")

	jwtToken, err := c.authService.AuthWithGoogleWorkspaceAccount(ctx.Request.Context(), code)
	if err != nil {
		if e, ok := err.(*errors.ErrInvalidGoogleWorkspaceAccount); ok {
			ctx.Redirect(http.StatusFound, redirect+fmt.Sprintf("&error=%v 로 끝나는 메일 주소만 사용 가능 합니다", e.Domain))
			return
		}

		ctx.Redirect(http.StatusFound, redirect+"&error=server-internal-error")
		return
	}

	refreshToken, err := ctx.Request.Cookie("refreshToken")
	if err != nil || len(refreshToken.Value) == 0 {
		cookie := new(http.Cookie)
		cookie.Name = "refreshToken"
		cookie.Value = jwtToken.RefreshToken
		cookie.HttpOnly = true
		cookie.Path = "/"
		cookie.Expires = jwtToken.GetRefreshTokenExpiresForCookie()

		http.SetCookie(ctx.Writer, cookie)
	} else {
		refreshToken.Value = jwtToken.RefreshToken
		refreshToken.HttpOnly = true
		refreshToken.Path = "/"
		refreshToken.Expires = jwtToken.GetRefreshTokenExpiresForCookie()

		http.SetCookie(ctx.Writer, refreshToken)
	}

	ctx.Redirect(http.StatusFound, redirect+"&accessToken="+jwtToken.AccessToken)
}

func (c AuthController) refreshAccessToken(ctx *gin.Context) {
	var request map[string]string

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtAuthentication := security.JwtAuthentication{}
	accessToken, expiresAt, err := jwtAuthentication.RefreshAccessToken(request["refreshToken"])

	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	err = c.logMemberAccessAtByToken(ctx.Request.Context(), request["refreshToken"])
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	result := map[string]any{}
	result["accessToken"] = accessToken
	result["expiresAt"] = expiresAt

	ctx.JSON(http.StatusOK, result)
}

func (c AuthController) logMemberAccessAtByToken(ctx context.Context, token string) error {
	jwtAuthentication := security.JwtAuthentication{}
	userClaim, err := jwtAuthentication.ConvertTokenUserClaim(token)
	if err != nil {
		return err
	}

	err = c.memberService.UpdateMemberLastAccessAt(ctx, userClaim.Id)
	if err != nil {
		return err
	}

	return nil
}
