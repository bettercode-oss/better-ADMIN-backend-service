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

func TestMemberAccessLogController_LogMemberAccess_API_접근(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"type": "API_ACCESS",
		"url": "http://localhost:2016/api/members",
		"method": "GET",
		"parameters": "{\"page\":1,\"types\":\"dooray\"}",
		"payload": "{\"name\":\"231312dsax\"}"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/member-access-logs", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberAccessLogController{}.LogMemberAccess, ctx)

	// then
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMemberAccessLogController_LogMemberAccess_화면_접근(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"type": "PAGE_ACCESS",
		"url": "http://localhost:3306/#/memberss"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/member-access-logs", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberAccessLogController{}.LogMemberAccess, ctx)

	// then
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMemberAccessLogController_LogMemberAccess_지원_하지_않는_TYPE(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"type": "AACC",
		"url": "http://localhost:3306/#/memberss"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/member-access-logs", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberAccessLogController{}.LogMemberAccess, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMemberAccessLogController_GetMemberAccessLogs(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/member-access-logs?page=1&pageSize=2", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberAccessLogController{}.GetMemberAccessLogs, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":               float64(3),
				"memberId":         float64(1),
				"type":             "PAGE_ACCESS",
				"typeName":         "화면",
				"url":              "http://localhost:3000/#/menus",
				"ipAddress":        "127.0.0.1",
				"browserUserAgent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
				"createdAt":        "1982-01-04T00:00:00Z",
			},
			map[string]interface{}{
				"id":               float64(2),
				"memberId":         float64(1),
				"type":             "API_ACCESS",
				"typeName":         "API",
				"url":              "http://localhost:2016/api/members",
				"method":           "GET",
				"parameters":       "{\"page\":1,\"types\":\"dooray\"}",
				"payload":          "{\"name\":\"231312dsax\"}",
				"ipAddress":        "127.0.0.1",
				"browserUserAgent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
				"createdAt":        "1982-01-04T00:00:00Z",
			},
		},
		"totalCount": float64(3),
	}

	assert.Equal(t, expected, resp)
}

func TestMemberAccessLogController_GetMemberAccessLogs_By_멤버_아이디(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/member-access-logs?page=1&pageSize=2&memberId=1", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberAccessLogController{}.GetMemberAccessLogs, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":               float64(3),
				"memberId":         float64(1),
				"type":             "PAGE_ACCESS",
				"typeName":         "화면",
				"url":              "http://localhost:3000/#/menus",
				"ipAddress":        "127.0.0.1",
				"browserUserAgent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
				"createdAt":        "1982-01-04T00:00:00Z",
			},
			map[string]interface{}{
				"id":               float64(2),
				"memberId":         float64(1),
				"type":             "API_ACCESS",
				"typeName":         "API",
				"url":              "http://localhost:2016/api/members",
				"method":           "GET",
				"parameters":       "{\"page\":1,\"types\":\"dooray\"}",
				"payload":          "{\"name\":\"231312dsax\"}",
				"ipAddress":        "127.0.0.1",
				"browserUserAgent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
				"createdAt":        "1982-01-04T00:00:00Z",
			},
		},
		"totalCount": float64(2),
	}

	assert.Equal(t, expected, resp)

}
