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

func TestAuthController_AuthWithSignIdPassword(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"id": "siteadm",
		"password": "123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AuthController{}.AuthWithSignIdPassword, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp.(map[string]interface{})["accessToken"])
}

func TestAuthController_AuthWithSignIdPassword_미_승인_사용자(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"id": "ymyoo3",
		"password": "123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AuthController{}.AuthWithSignIdPassword, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNotAcceptable, rec.Code)
}
