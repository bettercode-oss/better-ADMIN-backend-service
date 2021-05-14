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

func TestAccessControlController_CreatePermission(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.CreatePermission, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_CreatePermission_권한명이_이미_있는_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"name": "MANAGE_MEMBERS"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/permissions", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.CreatePermission, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "duplicated", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_GetPermissions(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions?page=2&pageSize=2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.GetPermissions, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(3), resp.(map[string]interface{})["totalCount"])

	permissions := resp.(map[string]interface{})["result"].([]interface{})
	index := 0
	assert.Equal(t, float64(3), permissions[index].(map[string]interface{})["id"])
	assert.Equal(t, "user-define", permissions[index].(map[string]interface{})["type"])
	assert.Equal(t, "ACCESS_STOCK", permissions[index].(map[string]interface{})["name"])
	assert.Equal(t, "재고 접근 권한", permissions[index].(map[string]interface{})["description"])
}

func TestAccessControlController_UpdatePermission(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	permissionId := "3"
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/:permissionId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues(permissionId)

	// when
	handleWithFilter(AccessControlController{}.UpdatePermission, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_UpdatePermission_사전_정의_유형(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	permissionId := "2"
	requestBody := `{
		"name": "PRODUCT-MANGED",
		"description": "상품 관리 권한"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/:permissionId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues(permissionId)

	// when
	handleWithFilter(AccessControlController{}.UpdatePermission, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_UpdatePermission_이미_기존에_존재하는_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	permissionId := "3"
	requestBody := `{
		"name": "MANAGE_MEMBERS",
		"description": "기존에 존재하는 권한명"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/permissions/:permissionId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues(permissionId)

	// when
	handleWithFilter(AccessControlController{}.UpdatePermission, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "duplicated", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_DeletePermission(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	permissionId := "3"
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/:permissionId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues(permissionId)

	// when
	handleWithFilter(AccessControlController{}.DeletePermission, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_DeletePermission_사전_정의_유형(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	permissionId := "2"
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/permissions/:permissionId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues(permissionId)

	// when
	handleWithFilter(AccessControlController{}.DeletePermission, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_CreateRole(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"name": "MD",
		"description": "MD 역할",
    "allowedPermissionIds": [2, 3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/access-control/roles", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.CreateRole, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_GetRoles(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.GetRoles, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(3), resp.(map[string]interface{})["totalCount"])

	roles := resp.(map[string]interface{})["result"].([]interface{})
	index := 0
	assert.Equal(t, float64(1), roles[index].(map[string]interface{})["id"])
	assert.Equal(t, "pre-define", roles[index].(map[string]interface{})["type"])
	assert.Equal(t, "사전정의", roles[index].(map[string]interface{})["typeName"])
	assert.Equal(t, "SYSTEM MANAGER", roles[index].(map[string]interface{})["name"])
	assert.Equal(t, "시스템 관리자", roles[index].(map[string]interface{})["description"])

	permissions := roles[index].(map[string]interface{})["permissions"].([]interface{})
	assert.Equal(t, 2, len(permissions))
	permissionIndex := 0
	assert.Equal(t, float64(1), permissions[permissionIndex].(map[string]interface{})["id"])
	assert.Equal(t, "MANAGE_SYSTEM_SETTINGS", permissions[permissionIndex].(map[string]interface{})["name"])
	permissionIndex++
	assert.Equal(t, float64(2), permissions[permissionIndex].(map[string]interface{})["id"])
	assert.Equal(t, "MANAGE_MEMBERS", permissions[permissionIndex].(map[string]interface{})["name"])

	index++
	assert.Equal(t, float64(2), roles[index].(map[string]interface{})["id"])
	assert.Equal(t, "pre-define", roles[index].(map[string]interface{})["type"])
	assert.Equal(t, "사전정의", roles[index].(map[string]interface{})["typeName"])
	assert.Equal(t, "MEMBER MANAGER", roles[index].(map[string]interface{})["name"])
	assert.Equal(t, "멤버 관리자", roles[index].(map[string]interface{})["description"])

	permissions = roles[index].(map[string]interface{})["permissions"].([]interface{})
	assert.Equal(t, 1, len(permissions))
	permissionIndex = 0
	assert.Equal(t, float64(2), permissions[permissionIndex].(map[string]interface{})["id"])
	assert.Equal(t, "MANAGE_MEMBERS", permissions[permissionIndex].(map[string]interface{})["name"])

	index++
	assert.Equal(t, float64(3), roles[index].(map[string]interface{})["id"])
	assert.Equal(t, "user-define", roles[index].(map[string]interface{})["type"])
	assert.Equal(t, "사용자정의", roles[index].(map[string]interface{})["typeName"])
	assert.Equal(t, "테스트 관리자", roles[index].(map[string]interface{})["name"])
	assert.Equal(t, "", roles[index].(map[string]interface{})["description"])

	permissions = roles[index].(map[string]interface{})["permissions"].([]interface{})
	assert.Equal(t, 1, len(permissions))
	permissionIndex = 0
	assert.Equal(t, float64(1), permissions[permissionIndex].(map[string]interface{})["id"])
	assert.Equal(t, "MANAGE_SYSTEM_SETTINGS", permissions[permissionIndex].(map[string]interface{})["name"])
}

func TestAccessControlController_DeleteRole(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	roleId := "3"
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/:roleId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues(roleId)

	// when
	handleWithFilter(AccessControlController{}.DeleteRole, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_DeleteRole_사전정의_유형(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	roleId := "2"
	req := httptest.NewRequest(http.MethodDelete, "/api/access-control/roles/:roleId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues(roleId)

	// when
	handleWithFilter(AccessControlController{}.DeleteRole, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	fmt.Println(rec.Body.String())

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}

func TestAccessControlController_UpdateRole(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	roleId := "3"
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/role/:roleId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues(roleId)

	// when
	handleWithFilter(AccessControlController{}.UpdateRole, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAccessControlController_UpdateRole_사전정의_유형(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	roleId := "2"
	requestBody := `{
		"name": "프로덕트 오너",
		"description": "프로덕트",
    "allowedPermissionIds": [1, 2, 3]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/access-control/role/:roleId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues(roleId)

	// when
	handleWithFilter(AccessControlController{}.UpdateRole, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}
