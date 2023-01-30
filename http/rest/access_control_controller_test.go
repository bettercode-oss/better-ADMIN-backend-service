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

func TestAccessControlController_createPermission_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_createPermission_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAccessControlController_createPermission(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_createPermission_권한명이_이미_있는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "MANAGE_MEMBERS"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, "duplicated", actual.(map[string]interface{})["message"])
}

func TestAccessControlController_getPermissions_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions?page=2&pageSize=2", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
}

func TestAccessControlController_getPermissions(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions?page=2&pageSize=2", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(3),
				"type":        "user-define",
				"typeName":    "사용자정의",
				"name":        "ACCESS_STOCK",
				"description": "재고 접근 권한",
			},
		},
		"totalCount": float64(3),
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getPermissions_이름으로_검색(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions?page=1&pageSize=10&name=ACCESS", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(3),
				"type":        "user-define",
				"typeName":    "사용자정의",
				"name":        "ACCESS_STOCK",
				"description": "재고 접근 권한",
			},
		},
		"totalCount": float64(1),
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getPermission_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
}

func TestAccessControlController_getPermission(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"id":          float64(3),
		"type":        "user-define",
		"typeName":    "사용자정의",
		"name":        "ACCESS_STOCK",
		"description": "재고 접근 권한",
		"createdAt":   "1982-01-04T00:00:00Z",
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getPermission_ID에_해당하는_권한이_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	fmt.Println(rec.Body.String())
}

func TestAccessControlController_updatePermission_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAccessControlController_updatePermission_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_updatePermission_permission_id가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/1000", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAccessControlController_updatePermission(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_UpdatePermission_사전_정의_유형(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/2", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, "non changeable", actual.(map[string]interface{})["message"])
}

func TestAccessControlController_UpdatePermission_이미_기존에_존재하는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "MANAGE_MEMBERS",
		"description": "기존에 존재하는 권한명"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, "duplicated", actual.(map[string]interface{})["message"])
}

func TestAccessControlController_deletePermission_권한_확인(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
}

func TestAccessControlController_deletePermission_member_id_가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_deletePermission(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_deletePermission_사전_정의_유형(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/2", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_createRole_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "MD",
		"description": "MD 역할",
    "allowedPermissionIds": [2, 3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAccessControlController_createRole_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "MD 역할",
    "allowedPermissionIds": [2, 3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_createRole(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "MD",
		"description": "MD 역할",
    "allowedPermissionIds": [2, 3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_getRoles_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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

func TestAccessControlController_getRoles(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(1),
				"type":        "pre-define",
				"typeName":    "사전정의",
				"name":        "SYSTEM MANAGER",
				"description": "시스템 관리자",
				"permissions": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "MANAGE_SYSTEM_SETTINGS",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "MANAGE_MEMBERS",
					},
				},
			},
			map[string]interface{}{
				"id":          float64(2),
				"type":        "pre-define",
				"typeName":    "사전정의",
				"name":        "MEMBER MANAGER",
				"description": "멤버 관리자",
				"permissions": []interface{}{
					map[string]interface{}{
						"id":   float64(2),
						"name": "MANAGE_MEMBERS",
					},
				},
			},
			map[string]interface{}{
				"id":          float64(3),
				"type":        "user-define",
				"typeName":    "사용자정의",
				"name":        "테스트 관리자",
				"description": "",
				"permissions": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "MANAGE_SYSTEM_SETTINGS",
					},
				},
			},
		},
		"totalCount": float64(3),
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getRoles_이름으로_검색(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles?name=테스", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(3),
				"type":        "user-define",
				"typeName":    "사용자정의",
				"name":        "테스트 관리자",
				"description": "",
				"permissions": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "MANAGE_SYSTEM_SETTINGS",
					},
				},
			},
		},
		"totalCount": float64(1),
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getRole_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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

func TestAccessControlController_getRole(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"id":          float64(3),
		"type":        "user-define",
		"typeName":    "사용자정의",
		"name":        "테스트 관리자",
		"description": "",
		"createdAt":   "1982-01-04T00:00:00Z",
		"permissions": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "MANAGE_SYSTEM_SETTINGS",
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestAccessControlController_getRole_ID가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	fmt.Println(rec.Body.String())
}

func TestAccessControlController_updateRole_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/roles/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAccessControlController_updateRole_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/roles/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_updateRole_role_id가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/roles/1000", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAccessControlController_updateRole(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/roles/3", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_updateRole_사전정의_유형(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
   "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/roles/2", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, "non changeable", actual.(map[string]interface{})["message"])
}

func TestAccessControlController_deleteRole_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"TC",
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
}

func TestAccessControlController_deleteRole_role_id_가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_deleteRole(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/3", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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

func TestAccessControlController_deleteRole_사전정의_유형(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/2", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ACCESS_CONTROL",
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
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	fmt.Println(rec.Body.String())

	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.Equal(t, "non changeable", actual.(map[string]interface{})["message"])
}
