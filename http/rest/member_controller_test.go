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

func TestMemberController_signUpMember(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"signId": "ymyoo1",
		"name": "유영모",
		"password": "1111"	
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/members", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMemberController_signUpMember_아이디_중복(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"signId": "ymyoo",
		"name": "유영모",
		"password": "1111"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/members", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMemberController_signUpMember_필수값_확인(t *testing.T) {
	// given
	requestBody := `{
		"name": "유영모",
		"password": "1111"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/members", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMemberController_getCurrentMember_토큰이_없는_경우(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/my", nil)
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMemberController_getCurrentMember_MemberId가_유효하지_않는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/my", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1000,
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
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestMemberController_getCurrentMember(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/my", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
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
		"id":          float64(1),
		"type":        "site",
		"typeName":    "사이트",
		"name":        "사이트 관리자",
		"roles":       []any{"SYSTEM MANAGER", "MEMBER MANAGER"},
		"permissions": []any{"MANAGE_SYSTEM_SETTINGS", "MANAGE_MEMBERS"},
		"picture":     "",
	}
	assert.Equal(t, expected, actual)
}

func TestMemberController_getMembers_by_멤버_역할(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&roleIds=1", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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
		"result": []any{
			map[string]any{
				"id":           float64(1),
				"signId":       "siteadm",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "siteadm",
				"name":         "사이트 관리자",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]any{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]any{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
		},
		"totalCount": float64(2),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberController_getMembers_승인된_멤버(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=2&status=approved", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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
		"result": []any{
			map[string]any{
				"id":          float64(1),
				"signId":      "siteadm",
				"candidateId": "siteadm",
				"type":        "site",
				"typeName":    "사이트",
				"name":        "사이트 관리자",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
			},
			map[string]any{
				"id":          float64(2),
				"signId":      "",
				"candidateId": "2222",
				"type":        "dooray",
				"typeName":    "두레이",
				"name":        "유영모",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]any{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
			},
		},
		"totalCount": float64(3),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberController_getMembers_신청한_멤버(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=applied", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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
		"result": []any{
			map[string]any{
				"id":            float64(4),
				"signId":        "ymyoo3",
				"candidateId":   "ymyoo3",
				"type":          "site",
				"typeName":      "사이트",
				"name":          "유영모3",
				"roles":         []any{},
				"organizations": []any{},
				"createdAt":     "1982-01-04T00:00:00Z",
				"lastAccessAt":  "1982-01-05T00:00:00Z",
			},
		},
		"totalCount": float64(1),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberController_getMembers_by_멤버_이름(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&name=유", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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
		"result": []any{
			map[string]any{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]any{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]any{
				"id":           float64(3),
				"signId":       "ymyoo",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "ymyoo",
				"name":         "유영모2",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles":        []any{},
				"organizations": []any{
					map[string]any{
						"id":   float64(4),
						"name": "부서C",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
						},
					},
				},
			},
		},
		"totalCount": float64(2),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberController_getMembers_by_멤버_유형(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&types=dooray,site", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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
		"result": []any{
			map[string]any{
				"id":           float64(1),
				"signId":       "siteadm",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "siteadm",
				"name":         "사이트 관리자",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]any{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []any{
					map[string]any{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]any{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []any{
					map[string]any{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]any{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]any{
				"id":           float64(3),
				"signId":       "ymyoo",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "ymyoo",
				"name":         "유영모2",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles":        []any{},
				"organizations": []any{
					map[string]any{
						"id":   float64(4),
						"name": "부서C",
						"roles": []any{
							map[string]any{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
						},
					},
				},
			},
		},
		"totalCount": float64(3),
	}

	assert.Equal(t, expected, actual)
}

func TestMemberController_getMember_권한이_없는_경우(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/1", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
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

func TestMemberController_getMember_member_id_가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/1000", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.read",
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

func TestMemberController_getMember(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/1", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.read",
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
	assert.Equal(t, float64(1), actual.(map[string]any)["id"])
	assert.Equal(t, "site", actual.(map[string]any)["type"])
	assert.Equal(t, "사이트", actual.(map[string]any)["typeName"])
	assert.Equal(t, "사이트 관리자", actual.(map[string]any)["name"])

	memberRoles := actual.(map[string]any)["roles"].([]any)
	assert.Equal(t, 1, len(memberRoles))
	memberRoleIndex := 0
	assert.Equal(t, float64(1), memberRoles[memberRoleIndex].(map[string]any)["id"])
	assert.Equal(t, "SYSTEM MANAGER", memberRoles[memberRoleIndex].(map[string]any)["name"])
}

func TestMemberController_assignRole_Bad_Request_필수_값_확인(t *testing.T) {
	// given
	requestBody := `{
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/1/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_assignRole_권한_확인(t *testing.T) {
	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/1/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.read",
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

func TestMemberController_assignRole_member_id가_유효하지_않는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/1000/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_assignRole(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/1/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_assignRoles_역할이_없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"roleIds": []
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/1/assign-roles", strings.NewReader(requestBody))
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_approveMember_member_id_가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/1000/approved", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_approveMember_권한_확인(t *testing.T) {
	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/4/approved", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.read",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestMemberController_approveMember(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/4/approved", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_approveMember_이미_승인된_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/1/approved", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_rejectMember_권한_확인(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/4/rejected", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.read",
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
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestMemberController_rejectMember_member_id_가_유효하지_않은_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/10000/rejected", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_rejectMember(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/4/rejected", nil)
	token, err := generateTestJWT(map[string]any{
		"Id": 1,
		"Permissions": []string{
			"member.update",
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

func TestMemberController_getSearchFilters(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/search-filters", nil)
	token, err := generateTestJWT(map[string]any{
		"Id":          1,
		"Permissions": []string{},
	}, time.Minute*15)

	if err != nil {
		t.Failed()
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// when
	ginApp.ServeHTTP(rec, req)

	fmt.Println(rec.Body.String())
	var actual any
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := []any{
		map[string]any{
			"name": "type",
			"filters": []any{
				map[string]any{
					"text":  "사이트",
					"value": "site",
				},
				map[string]any{
					"text":  "두레이",
					"value": "dooray",
				},
				map[string]any{
					"text":  "구글",
					"value": "google",
				},
			},
		},
		map[string]any{
			"name": "role",
			"filters": []any{
				map[string]any{
					"text":  "SYSTEM MANAGER",
					"value": "1",
				},
				map[string]any{
					"text":  "MEMBER MANAGER",
					"value": "2",
				},
				map[string]any{
					"text":  "테스트 관리자",
					"value": "3",
				},
			},
		},
	}

	assert.Equal(t, expected, actual)
}
