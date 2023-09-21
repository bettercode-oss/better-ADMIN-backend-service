package rest

import (
	"better-admin-backend-service/config"
	"better-admin-backend-service/security"
	"better-admin-backend-service/testdata/testdb"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_authWithSignIdPassword(t *testing.T) {
	// setup Fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
    "id": "siteadm",
    "password": "123456"
  }`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
	var actual map[string]any
	json.Unmarshal(rec.Body.Bytes(), &actual)
	accessTokenClaim, err := security.JwtAuthentication{}.ConvertTokenUserClaim(actual["accessToken"].(string))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, accessTokenClaim.Id, uint(1))
	assert.Equal(t, accessTokenClaim.Roles, []string{"SYSTEM MANAGER", "MEMBER MANAGER"})
	assert.Equal(t, accessTokenClaim.Permissions, []string{"MANAGE_SYSTEM_SETTINGS", "MANAGE_MEMBERS"})

	refreshTokenClaim, err := security.JwtAuthentication{}.ConvertTokenUserClaim(actual["refreshToken"].(string))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, refreshTokenClaim.Id, uint(1))
	assert.Equal(t, refreshTokenClaim.Roles, []string{"SYSTEM MANAGER", "MEMBER MANAGER"})
	assert.Equal(t, refreshTokenClaim.Permissions, []string{"MANAGE_SYSTEM_SETTINGS", "MANAGE_MEMBERS"})
}

func Test_authWithSignIdPassword_Bad_Request(t *testing.T) {
	// given
	requestBody := `{
    "id": "siteadm"
  }`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_authWithSignIdPassword_계정이_유효하지_않는_경우(t *testing.T) {
	// setup Fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
    "id": "Got7",
		"password": "123456"
  }`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_authWithSignIdPassword_비밀번호가_유효하지_않은_경우(t *testing.T) {
	// setup Fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
    "id": "siteadm",
		"password": "qwert"
  }`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func Test_authWithSignIdPassword_미_승인_사용자(t *testing.T) {
	// setup Fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"id": "ymyoo3",
		"password": "123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNotAcceptable, rec.Code)
}

func Test_authWithGoogleWorkspaceAccount(t *testing.T) {
	// setup fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// Google Workspace Server Fixture
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

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.True(t, strings.Contains(rec.Header().Get("Location"), "accessToken="))

	// assert Cookie value
	headerSetCookie := rec.Header().Get("Set-Cookie")
	fmt.Println("Set-Cookie in headers", headerSetCookie)

	assert.NotEmpty(t, headerSetCookie)
	assert.True(t, strings.HasPrefix(headerSetCookie, "refreshToken="))
	expires := "Expires=" + time.Now().Add(time.Hour*24*7).Format("Mon, 02 Jan 2006")
	assert.True(t, strings.Contains(headerSetCookie, expires))
	assert.True(t, strings.Contains(headerSetCookie, "HttpOnly"))

	refreshToken := headerSetCookie[strings.Index(headerSetCookie, "refreshToken=")+len("refreshToken=") : strings.Index(headerSetCookie, ";")]
	tokenUserClaim, _ := security.JwtAuthentication{}.ConvertTokenUserClaim(refreshToken)
	assert.Equal(t, uint(5), tokenUserClaim.Id)
}

func Test_authWithGoogleWorkspaceAccount_구글_워크스페이스_멤버가_아닌_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

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

	googleWorkspaceServerUrl := fmt.Sprintf("http://localhost:%v", serverPort)
	config.Config.GoogleOAuth.AuthUri = googleWorkspaceServerUrl
	config.Config.GoogleOAuth.TokenUri = googleWorkspaceServerUrl

	// given
	code := "test-google-code"
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/auth/google-workspace?code=%v", code), nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println("Location", rec.Header().Get("Location"))
	fmt.Println("Set-Cookie", rec.Header().Get("Set-Cookie"))
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.False(t, strings.Contains(rec.Header().Get("Location"), "accessToken="))
	assert.False(t, strings.Contains(rec.Header().Get("Set-Cookie"), "refreshToken="))
	// 받환 메시지( error=bettercode.kr 로 끝나는 메일 주소만 사용 가능 합니다.) 중 한글 부분이 인코딩되어 있기 때문에 인코딩 값을 비교
	assert.True(t, strings.Contains(rec.Header().Get("Location"),
		"error=bettercode.kr %eb%a1%9c %eb%81%9d%eb%82%98%eb%8a%94 %eb%a9%94%ec%9d%bc %ec%a3%bc%ec%86%8c%eb%a7%8c %ec%82%ac%ec%9a%a9 %ea%b0%80%eb%8a%a5 %ed%95%a9%eb%8b%88%eb%8b%a4"))
}

func Test_refreshAccessToken(t *testing.T) {
	// setup Fixture
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	refreshToken, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Roles":       []string{},
		"Permissions": []string{},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}

	requestBody := fmt.Sprintf(`{
		"refreshToken": "%s"
	}`, refreshToken)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/token/refresh", strings.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NotEmpty(t, actual.(map[string]interface{})["accessToken"])
}

func Test_refreshAccessToken_토큰이_없는_경우(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodPost, "/api/auth/token/refresh", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
