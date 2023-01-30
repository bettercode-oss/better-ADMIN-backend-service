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

func TestOrganizationController_createOrganization_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "테스트 조직"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
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

func TestOrganizationController_createOrganization_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_createOrganization_최상위_조직으로_추가(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "테스트 조직"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_createOrganization_상위조직이_있는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"parentOrganizationId": 1,
		"name": "테스트 조직"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_getOrganizations_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations", nil)
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

func TestOrganizationController_getOrganizations(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

	expected := []interface{}{
		map[string]interface{}{
			"id":   float64(1),
			"name": "베터코드 연구소",
			"subOrganizations": []interface{}{
				map[string]interface{}{
					"id":   float64(3),
					"name": "부서B",
					"subOrganizations": []interface{}{
						map[string]interface{}{
							"id":   float64(4),
							"name": "부서C",
							"roles": []interface{}{
								map[string]interface{}{
									"id":   float64(1),
									"name": "SYSTEM MANAGER",
								},
							},
							"members": []interface{}{
								map[string]interface{}{
									"id":   float64(3),
									"name": "유영모2",
								},
							},
						},
					},
				},
			},
			"roles": []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "SYSTEM MANAGER",
				}, map[string]interface{}{
					"id":   float64(2),
					"name": "MEMBER MANAGER",
				},
			},
			"members": []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "사이트 관리자",
				}, map[string]interface{}{
					"id":   float64(2),
					"name": "유영모",
				},
			},
		},
		map[string]interface{}{
			"id":   float64(5),
			"name": "베터코드 연구소2",
			"subOrganizations": []interface{}{
				map[string]interface{}{
					"id":   float64(2),
					"name": "부서A",
				},
			},
		},
	}
	assert.Equal(t, expected, actual.([]interface{}))
}

func TestOrganizationController_getOrganization_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/1", nil)
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

func TestOrganizationController_getOrganization(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/1", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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
		"id":        float64(1),
		"name":      "베터코드 연구소",
		"createdAt": "1982-01-04T00:00:00Z",
		"roles": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "SYSTEM MANAGER",
			}, map[string]interface{}{
				"id":   float64(2),
				"name": "MEMBER MANAGER",
			},
		},
		"members": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "사이트 관리자",
			}, map[string]interface{}{
				"id":   float64(2),
				"name": "유영모",
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestOrganizationController_getOrganization_ID_로_찾을수없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_changeOrganizationName_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "강남 베터코드"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1/name", strings.NewReader(requestBody))
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

func TestOrganizationController_changeOrganizationName_id_가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "강남 베터코드"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1000/name", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_changeOrganizationName(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "강남 베터코드"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1/name", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_changePosition_id_가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"parentOrganizationId": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1000/change-position", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_changePosition_하위로_변경(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"parentOrganizationId": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/2/change-position", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_changePosition_최상위로_변경(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/2/change-position", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_assignRoles_id_가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1000/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_AssignRoles(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_assignMembers_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1000/assign-members", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_assignMembers_id_가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"memberIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1000/assign-members", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_assignMembers(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"memberIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/1/assign-members", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_deleteOrganization_id가_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/organizations/1000", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_DeleteOrganization_최하위(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/organizations/4", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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

func TestOrganizationController_DeleteOrganization_최상위(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodDelete, "/api/organizations/1", nil)
	token, err := generateTestJWT(map[string]interface{}{
		"Id": 1,
		"Permissions": []string{
			"MANAGE_ORGANIZATION",
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
