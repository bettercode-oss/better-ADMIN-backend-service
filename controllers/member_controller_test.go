package controllers

import (
	"better-admin-backend-service/security"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMemberController_GetMembers_승인된_멤버(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=2&status=approved", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":          float64(1),
				"signId":      "siteadm",
				"candidateId": "siteadm",
				"type":        "site",
				"typeName":    "사이트",
				"name":        "사이트 관리자",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
			},
			map[string]interface{}{
				"id":          float64(2),
				"signId":      "",
				"candidateId": "2222",
				"type":        "dooray",
				"typeName":    "두레이",
				"name":        "유영모",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
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

	assert.Equal(t, expected, resp)

}

func TestMemberController_GetMembers_신청한_멤버(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=applied", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":            float64(4),
				"signId":        "ymyoo3",
				"candidateId":   "ymyoo3",
				"type":          "site",
				"typeName":      "사이트",
				"name":          "유영모3",
				"roles":         []interface{}{},
				"organizations": []interface{}{},
				"createdAt":     "1982-01-04T00:00:00Z",
				"lastAccessAt":  "1982-01-05T00:00:00Z",
			},
		},
		"totalCount": float64(1),
	}

	assert.Equal(t, expected, resp)
}

func TestMemberController_GetMembers_by_멤버_이름(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&name=유", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]interface{}{
				"id":           float64(3),
				"signId":       "ymyoo",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "ymyoo",
				"name":         "유영모2",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles":        []interface{}{},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(4),
						"name": "부서C",
						"roles": []interface{}{
							map[string]interface{}{
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

	assert.Equal(t, expected, resp)
}

func TestMemberController_GetMembers_by_멤버_유형(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&types=dooray,site", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":           float64(1),
				"signId":       "siteadm",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "siteadm",
				"name":         "사이트 관리자",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]interface{}{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]interface{}{
				"id":           float64(3),
				"signId":       "ymyoo",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "ymyoo",
				"name":         "유영모2",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles":        []interface{}{},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(4),
						"name": "부서C",
						"roles": []interface{}{
							map[string]interface{}{
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

	assert.Equal(t, expected, resp)
}

func TestMemberController_GetMembers_by_멤버_역할(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=10&status=approved&roleIds=1", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := map[string]interface{}{
		"result": []interface{}{
			map[string]interface{}{
				"id":           float64(1),
				"signId":       "siteadm",
				"type":         "site",
				"typeName":     "사이트",
				"candidateId":  "siteadm",
				"name":         "사이트 관리자",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
								"id":   float64(2),
								"name": "MEMBER MANAGER",
							},
						},
					},
				},
			},
			map[string]interface{}{
				"id":           float64(2),
				"signId":       "",
				"type":         "dooray",
				"typeName":     "두레이",
				"candidateId":  "2222",
				"name":         "유영모",
				"createdAt":    "1982-01-04T00:00:00Z",
				"lastAccessAt": "1982-01-05T00:00:00Z",
				"roles": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "SYSTEM MANAGER",
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "MEMBER MANAGER",
					},
				},
				"organizations": []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "베터코드 연구소",
						"roles": []interface{}{
							map[string]interface{}{
								"id":   float64(1),
								"name": "SYSTEM MANAGER",
							},
							map[string]interface{}{
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

	assert.Equal(t, expected, resp)
}

func TestMemberController_AssignRoles(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/:id/assign-roles", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberController{}.AssignRole, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMemberController_AssignRoles_역할이_없는_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"roleIds": []
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/members/:id/assign-roles", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberController{}.AssignRole, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMemberController_GetMember(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	// when
	handleWithFilter(MemberController{}.GetMember, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(1), resp.(map[string]interface{})["id"])
	assert.Equal(t, "site", resp.(map[string]interface{})["type"])
	assert.Equal(t, "사이트", resp.(map[string]interface{})["typeName"])
	assert.Equal(t, "사이트 관리자", resp.(map[string]interface{})["name"])

	memberRoles := resp.(map[string]interface{})["roles"].([]interface{})
	assert.Equal(t, 1, len(memberRoles))
	memberRoleIndex := 0
	assert.Equal(t, float64(1), memberRoles[memberRoleIndex].(map[string]interface{})["id"])
	assert.Equal(t, "SYSTEM MANAGER", memberRoles[memberRoleIndex].(map[string]interface{})["name"])
}

func TestMemberController_SignUpMember(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"signId": "ymyoo1",
		"name": "유영모",
		"password": "1111"	
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/members", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.SignUpMember, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMemberController_SignUpMember_아이디_중복(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"signId": "ymyoo",
		"name": "유영모",
		"password": "1111"	
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/members", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.SignUpMember, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMemberController_ApproveMember(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/:id/approved", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("4")

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberController{}.ApproveMember, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMemberController_ApproveMember_이미_승인된_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/:id/approved", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberController{}.ApproveMember, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMemberController_GetSearchFilters(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members/search-filters", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetSearchFilters, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := []interface{}{
		map[string]interface{}{
			"name": "type",
			"filters": []interface{}{
				map[string]interface{}{
					"text":  "사이트",
					"value": "site",
				},
				map[string]interface{}{
					"text":  "두레이",
					"value": "dooray",
				},
				map[string]interface{}{
					"text":  "구글",
					"value": "google",
				},
			},
		},
		map[string]interface{}{
			"name": "role",
			"filters": []interface{}{
				map[string]interface{}{
					"text":  "SYSTEM MANAGER",
					"value": "1",
				},
				map[string]interface{}{
					"text":  "MEMBER MANAGER",
					"value": "2",
				},
				map[string]interface{}{
					"text":  "테스트 관리자",
					"value": "3",
				},
			},
		},
	}

	assert.Equal(t, expected, resp)
}

func TestMemberController_RejectMember(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodPut, "/api/members/:id/rejected", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("4")

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MemberController{}.RejectMember, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
