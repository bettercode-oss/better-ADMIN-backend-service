package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"better-admin-backend-service/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AuthController struct {
}

func (controller AuthController) Init(rg *gin.RouterGroup) {
	route := rg.Group("/auth")

	route.POST("", controller.authWithSignIdPassword)
	route.POST("/dooray", controller.authWithDoorayIdPassword)
	route.GET("/google-workspace", controller.authWithGoogleWorkspaceAccount)
	route.GET("/check", controller.checkAuth)
	route.POST("/logout", controller.logout)
	route.POST("/token/refresh", controller.refreshAccessToken)
}

func (AuthController) authWithSignIdPassword(ctx *gin.Context) {
	var memberSignIn dtos.MemberSignIn

	if err := ctx.BindJSON(&memberSignIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := services.AuthService{}.AuthWithSignIdPassword(ctx.Request.Context(), memberSignIn)
	if err != nil {
		if err == domain.ErrNotFound || err == domain.ErrAuthentication {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err == domain.ErrUnApproved {
			ctx.JSON(http.StatusNotAcceptable, err.Error())
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
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

	result := map[string]string{}
	result["accessToken"] = jwtToken.AccessToken

	ctx.JSON(http.StatusOK, result)
}

func (AuthController) authWithDoorayIdPassword(ctx *gin.Context) {
	var memberSignIn dtos.MemberSignIn

	if err := ctx.BindJSON(&memberSignIn); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jwtToken, err := services.AuthService{}.AuthWithDoorayIdAndPassword(ctx.Request.Context(), memberSignIn)
	if err != nil {
		if err == domain.ErrAuthentication {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		helpers.ErrorHelper().InternalServerError(ctx, err)
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

	result := map[string]string{}
	result["accessToken"] = jwtToken.AccessToken

	ctx.JSON(http.StatusOK, result)
}

func (AuthController) authWithGoogleWorkspaceAccount(ctx *gin.Context) {
	code := ctx.Query("code")
	redirect := ctx.Query("state")

	jwtToken, err := services.AuthService{}.AuthWithGoogleWorkspaceAccount(ctx.Request.Context(), code)
	if err != nil {
		if e, ok := err.(*domain.ErrInvalidGoogleWorkspaceAccount); ok {
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

func (AuthController) checkAuth(ctx *gin.Context) {
	refreshToken, err := ctx.Request.Cookie("refreshToken")
	if err != nil || len(refreshToken.Value) == 0 {
		ctx.JSON(http.StatusNotAcceptable, nil)
		return
	}

	jwtAuthentication := security.JwtAuthentication{}
	if err := jwtAuthentication.ValidateToken(refreshToken.Value); err != nil {
		log.Error(err)
		ctx.JSON(http.StatusNotAcceptable, nil)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (AuthController) logout(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	http.SetCookie(ctx.Writer, cookie)

	ctx.Status(http.StatusNoContent)
}

func (controller AuthController) refreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	refreshToken := cookie.Value
	jwtAuthentication := security.JwtAuthentication{}
	accessToken, err := jwtAuthentication.RefreshAccessToken(refreshToken)

	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	err = controller.logMemberAccessAtByToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		helpers.ErrorHelper().InternalServerError(ctx, err)
		return
	}

	result := map[string]string{}
	result["accessToken"] = accessToken
	ctx.JSON(http.StatusOK, result)
}

func (AuthController) logMemberAccessAtByToken(ctx context.Context, token string) error {
	jwtAuthentication := security.JwtAuthentication{}
	userClaim, err := jwtAuthentication.ConvertTokenUserClaim(token)
	if err != nil {
		return err
	}

	err = services.MemberService{}.UpdateMemberLastAccessAt(ctx, userClaim.Id)
	if err != nil {
		return err
	}

	return nil
}
