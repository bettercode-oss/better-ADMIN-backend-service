package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoRoute_Route에_등록된_URL과_Method는_통과_시킨다(t *testing.T) {
	// given
	router := gin.Default()
	router.Use(NoRoute(router))
	router.GET("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	router.POST("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	router.ServeHTTP(rec, req)
	// then
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestNoRoute_Route에_등록된_URL과_다른_Method도_통과_시킨다(t *testing.T) {
	// given
	router := gin.Default()
	router.Use(NoRoute(router))
	router.GET("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	router.POST("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNoRoute_Route에_등록되지_않는_URL을_호출하면_HTTP_Status_Code를_404_NotFound_를_반환한다(t *testing.T) {
	// given
	router := gin.Default()

	router.Use(NoRoute(router))
	router.GET("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	router.POST("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodGet, "/api/not-found", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestNoRoute_Route에_URL은_등록되어있지만_등록되지_않는_Method_호출하면_HTTP_Status_Code를_404_NotFound_를_반환한다(t *testing.T) {
	// given
	router := gin.Default()
	router.Use(NoRoute(router))
	router.GET("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	router.POST("/api/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	// when
	req := httptest.NewRequest(http.MethodDelete, "/api/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
