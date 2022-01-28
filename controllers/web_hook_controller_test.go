package controllers

import (
	"better-admin-backend-service/security"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebHookController_CreateWebHook(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"name": "테스트 웹훅",
		"description": "설명...."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 1,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(WebHookController{}.CreateWebHook, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestWebHookController_GetWebHooks(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/web-hooks?page=1&pageSize=2", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(WebHookController{}.GetWebHooks, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"totalCount": float64(3),
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(1),
				"name":        "테스트 웹훅",
				"description": "...",
			},
			map[string]interface{}{
				"id":          float64(2),
				"name":        "테스트 웹훅2",
				"description": "...",
			},
		},
	}

	assert.Equal(t, expected, resp.(map[string]interface{}))
}

func TestWebHookController_DeleteWebHook(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	webHookId := "3"
	req := httptest.NewRequest(http.MethodDelete, "/api/web-hooks/:id", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(webHookId)

	userClaim := security.UserClaim{
		Id: 1,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(WebHookController{}.DeleteWebHook, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestWebHookController_GetWebHook(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	webHookId := "3"
	req := httptest.NewRequest(http.MethodGet, "/api/web-hooks/:id", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(webHookId)

	// when
	handleWithFilter(WebHookController{}.GetWebHook, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"id":          float64(3),
		"name":        "테스트 웹훅3",
		"description": "...",
		"webHookCallSpec": map[string]interface{}{
			"httpRequestMethod": "POST",
			"url":               "http://example.com/api/web-hooks/3/note",
			"accessToken":       "test-access-tokens3",
			"sampleRequest":     "curl -X POST http://example.com/api/web-hooks/3/note -H \"Content-Type: application/json\" -H \"Authorization: Bearer test-access-tokens3\" -d '{\"text\":\"테스트 메시지 입니다.\"}'",
		},
	}

	assert.Equal(t, expected, resp.(map[string]interface{}))
}

func TestWebHookController_UpdateWebHook(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	webHookId := "3"
	requestBody := `{
		"name": "테스트 웹훅45444",
		"description": "변경된 설명...."
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/web-hooks/:id", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(webHookId)

	userClaim := security.UserClaim{
		Id: 1,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(WebHookController{}.UpdateWebHook, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestWebHookController_NoteMessage(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	webHookId := "3"
	requestBody := `{
		"text": "테스트 메시지..."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks/:id/note", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(webHookId)

	// when
	handleWithFilter(WebHookController{}.NoteMessage, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
