package controllers

import (
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
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
