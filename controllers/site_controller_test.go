package controllers

import (
	"better-admin-backend-service/config"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSiteController_SetDoorayLoginSetting(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode",
		"authorizationToken": "test-token"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/dooray-login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(SiteController{}.SetDoorayLoginSetting, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestSiteController_GetDoorayLoginSetting(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/dooray-login", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(SiteController{}.GetDoorayLoginSetting, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, true, resp.(map[string]interface{})["used"])
	assert.Equal(t, "bettercode", resp.(map[string]interface{})["domain"])
	assert.Equal(t, "test token....", resp.(map[string]interface{})["authorizationToken"])
}

func TestSiteController_GetGoogleWorkspaceLoginSetting(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/google-workspace-login", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(SiteController{}.GetGoogleWorkspaceLoginSetting, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"used":         true,
		"domain":       "bettercode.kr",
		"clientId":     "test-client-id",
		"clientSecret": "test-secret",
		"redirectUri":  "http://localhost:2016",
	}

	assert.Equal(t, expected, actual)
}

func TestSiteController_SetGoogleWorkspaceLoginSetting(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode.kr",
		"clientId": "test-client-id",
		"clientSecret": "test-secret",
		"redirectUri": "http://localhost:2016"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/google-workspace-login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(SiteController{}.SetGoogleWorkspaceLoginSetting, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestSiteController_GetSettingsSummary(t *testing.T) {
	DatabaseFixture{}.setUpDefault()
	config.InitConfig("../config/config.json")

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(SiteController{}.GetSettingsSummary, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"doorayLoginUsed":          true,
		"googleWorkspaceLoginUsed": true,
		"googleWorkspaceOAuthUri":  "https://accounts.google.com/o/oauth2/auth?client_id=test-client-id&redirect_uri=http://localhost:2016&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&approval_prompt=force&access_type=offline",
	}

	assert.Equal(t, expected, actual)
}
