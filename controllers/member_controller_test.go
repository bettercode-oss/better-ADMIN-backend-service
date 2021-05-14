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

func TestMemberController_GetMembers(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/members?page=1&pageSize=2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MemberController{}.GetMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(3), resp.(map[string]interface{})["totalCount"])

	members := resp.(map[string]interface{})["result"].([]interface{})
	assert.Equal(t, 2, len(members))

	index := 0
	assert.Equal(t, float64(1), members[index].(map[string]interface{})["id"])
	assert.Equal(t, "site", members[index].(map[string]interface{})["type"])
	assert.Equal(t, "사이트", members[index].(map[string]interface{})["typeName"])
	assert.Equal(t, "사이트 관리자", members[index].(map[string]interface{})["name"])

	memberRoles := members[index].(map[string]interface{})["roles"].([]interface{})
	assert.Equal(t, 1, len(memberRoles))
	memberRoleIndex := 0
	assert.Equal(t, float64(1), memberRoles[memberRoleIndex].(map[string]interface{})["id"])
	assert.Equal(t, "SYSTEM MANAGER", memberRoles[memberRoleIndex].(map[string]interface{})["name"])

	index++
	assert.Equal(t, float64(2), members[index].(map[string]interface{})["id"])
	assert.Equal(t, "dooray", members[index].(map[string]interface{})["type"])
	assert.Equal(t, "두레이", members[index].(map[string]interface{})["typeName"])
	assert.Equal(t, "유영모", members[index].(map[string]interface{})["name"])

	memberRoles = members[index].(map[string]interface{})["roles"].([]interface{})
	assert.Equal(t, 2, len(memberRoles))
	memberRoleIndex = 0
	assert.Equal(t, float64(1), memberRoles[memberRoleIndex].(map[string]interface{})["id"])
	assert.Equal(t, "SYSTEM MANAGER", memberRoles[memberRoleIndex].(map[string]interface{})["name"])
	memberRoleIndex++
	assert.Equal(t, float64(2), memberRoles[memberRoleIndex].(map[string]interface{})["id"])
	assert.Equal(t, "MEMBER MANAGER", memberRoles[memberRoleIndex].(map[string]interface{})["name"])
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

	// when
	handleWithFilter(MemberController{}.AssignRole, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
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

	// when
	handleWithFilter(MemberController{}.AssignRole, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
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
