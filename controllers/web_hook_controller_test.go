package controllers

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

func TestWebHookController_createWebHook_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "설명...."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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

func TestWebHookController_CreateWebHook(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "테스트 웹훅",
		"description": "설명...."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestWebHookController_getWebHooks(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/web-hooks?page=1&pageSize=2", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	fmt.Println(rec.Body.String())
	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

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

	assert.Equal(t, expected, actual.(map[string]interface{}))
}

func TestWebHookController_getWebHook_id가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/web-hooks/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebHookController_getWebHook(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/web-hooks/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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

func TestWebHookController_deleteWebHook_id가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/web-hooks/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebHookController_DeleteWebHook(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/web-hooks/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestWebHookController_updateWebHook_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "변경된 설명...."
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/web-hooks/1000", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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

func TestWebHookController_updateWebHook_id가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "테스트 웹훅45444",
		"description": "변경된 설명...."
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/web-hooks/1000", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebHookController_updateWebHook(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "테스트 웹훅45444",
		"description": "변경된 설명...."
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/web-hooks/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_SYSTEM_SETTINGS",
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
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestWebHookController_noteMessage_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks/3/note", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"NOTE_WEB_HOOKS",
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

func TestWebHookController_noteMessage_id_가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"text": "테스트 메시지..."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks/1000/note", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"NOTE_WEB_HOOKS",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWebHookController_noteMessage(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"text": "테스트 메시지..."
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/web-hooks/3/note", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"NOTE_WEB_HOOKS",
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
	assert.Equal(t, http.StatusCreated, rec.Code)
}
