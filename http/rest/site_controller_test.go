package rest

import (
	"better-admin-backend-service/testdata/testdb"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSiteController_getSettingsSummary(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"doorayLoginUsed":          true,
		"googleWorkspaceLoginUsed": true,
		"googleWorkspaceOAuthUri":  "https://accounts.google.com/o/oauth2/auth?client_id=test-client-id&redirect_uri=http://localhost:2016&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&approval_prompt=force&access_type=offline",
	}

	assert.Equal(t, expected, actual)
}

func TestSiteController_getDoorayLoginSetting(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/dooray-login", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.read",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, true, actual.(map[string]any)["used"])
	assert.Equal(t, "bettercode", actual.(map[string]any)["domain"])
	assert.Equal(t, "test token....", actual.(map[string]any)["authorizationToken"])
}

func TestSiteController_getDoorayLoginSetting_토큰이_없는_경우(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/dooray-login", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	fmt.Println(rec.Body.String())
}

func TestSiteController_getDoorayLoginSetting_권한이_없는_경우(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/dooray-login", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"BTS",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusForbidden, rec.Code)
	fmt.Println(rec.Body.String())
}

func TestSiteController_setDoorayLoginSetting_Bad_Request_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"domain": "bettercode",
		"authorizationToken": "test-token"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/dooray-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSiteController_setDoorayLoginSetting_Bad_Request_used_true_일_때_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/dooray-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSiteController_setDoorayLoginSetting(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode",
		"authorizationToken": "test-token"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/dooray-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestSiteController_getGoogleWorkspaceLoginSetting(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/google-workspace-login", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.read",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"used":         true,
		"domain":       "bettercode.kr",
		"clientId":     "test-client-id",
		"clientSecret": "test-secret",
		"redirectUri":  "http://localhost:2016",
	}

	assert.Equal(t, expected, actual)
}

func TestSiteController_setGoogleWorkspaceLoginSetting(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode.kr",
		"clientId": "test-client-id",
		"clientSecret": "test-secret",
		"redirectUri": "http://localhost:2016"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/google-workspace-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestSiteController_setGoogleWorkspaceLoginSetting_Bad_Request_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
		"domain": "bettercode.kr",
		"clientId": "test-client-id",
		"clientSecret": "test-secret",
		"redirectUri": "http://localhost:2016"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/google-workspace-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSiteController_setGoogleWorkspaceLoginSetting_Bad_Request_used_가_true_일_때_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
		"used": true,
		"domain": "bettercode.kr",
		"redirectUri": "http://localhost:2016"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/google-workspace-login", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id":    1,
		"Roles": []string{},
		"Permissions": []string{
			"site-settings.update",
		},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	fmt.Println(rec.Body.String())
}

func TestSiteController_getAppVersion(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/site/settings/app-version", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println(rec.Body.String())

	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]any{
		"version": float64(2),
	}

	assert.Equal(t, expected, actual)
}

func TestSiteController_increaseAppVersion(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/site/settings/app-version", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
