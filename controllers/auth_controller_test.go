package controllers

import (
	"better-admin-backend-service/config"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net"
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

func TestAuthController_AuthWithGoogleWorkspaceAccount(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
					"id": "123456",
					"email": "gigamadness@gmail.com",
					"verified_email": true,
					"name": "유영모",
					"given_name": "영모",
					"family_name": "유",
					"hd": "bettercode.kr",
					"picture": "https://lh3.googleusercontent.com/a-/AOh14GgO6suMzX-rWsUVXV5cQZWVSdmWCSdDGG-9_LrqRQ=s96-c"
			}`
			w.Write([]byte(responseBody))
		} else if r.Method == http.MethodPost {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
				"access_token": "test-token",
				"expires_in": 3599,
				"refresh_token": "test-refresh-token",
				"scope": "https://www.googleapis.com/auth/userinfo.email openid",
				"token_type": "Bearer",
				"id_token": "test-id-token"
			}`
			w.Write([]byte(responseBody))

		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	url := fmt.Sprintf("http://localhost:%v", serverPort)
	config.Config.GoogleOAuth.AuthUri = url
	config.Config.GoogleOAuth.TokenUri = url

	// given
	code := "test-google-code"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/auth/google-workspace?code=%v", code), nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AuthController{}.AuthWithGoogleWorkspaceAccount, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.True(t, strings.Contains(rec.Header().Get("Location"), "accessToken="))
	assert.True(t, strings.Contains(rec.Header().Get("Set-Cookie"), "refreshToken="))
}

func TestAuthController_AuthWithGoogleWorkspaceAccount_구글_워크스페이스_멤버가_아닌_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
					"id": "123456",
					"email": "gigamadness@gmail.com",
					"verified_email": true,
					"name": "유영모",
					"given_name": "영모",
					"family_name": "유",
					"picture": "https://lh3.googleusercontent.com/a-/AOh14GgO6suMzX-rWsUVXV5cQZWVSdmWCSdDGG-9_LrqRQ=s96-c"
			}`
			w.Write([]byte(responseBody))
		} else if r.Method == http.MethodPost {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			responseBody := `{
				"access_token": "test-token",
				"expires_in": 3599,
				"refresh_token": "test-refresh-token",
				"scope": "https://www.googleapis.com/auth/userinfo.email openid",
				"token_type": "Bearer",
				"id_token": "test-id-token"
			}`
			w.Write([]byte(responseBody))

		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverPort := server.Listener.Addr().(*net.TCPAddr).Port

	url := fmt.Sprintf("http://localhost:%v", serverPort)
	config.Config.GoogleOAuth.AuthUri = url
	config.Config.GoogleOAuth.TokenUri = url

	// given
	code := "test-google-code"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/auth/google-workspace?code=%v", code), nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AuthController{}.AuthWithGoogleWorkspaceAccount, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.False(t, strings.Contains(rec.Header().Get("Location"), "accessToken="))
	assert.False(t, strings.Contains(rec.Header().Get("Set-Cookie"), "refreshToken="))
	assert.True(t, strings.Contains(rec.Header().Get("Location"), "error=bettercode.kr 로 끝나는 메일 주소만 사용 가능 합니다."))
}
