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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 1,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

func TestAccessControlController_GetPermissions_이름으로_검색(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permissions?page=1&pageSize=10&name=ACCESS", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(AccessControlController{}.GetPermissions, ctx)

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

func TestAccessControlController_GetPermission(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permission/:permissionId", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues("3")

	// when
	handleWithFilter(AccessControlController{}.GetPermission, ctx)

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

func TestAccessControlController_GetPermission_ID에_해당하는_권한이_없는_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/permission/:permissionId", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("permissionId")
	ctx.SetParamValues("100")

	// when
	handleWithFilter(AccessControlController{}.GetPermission, ctx)

	// then
	assert.Equal(t, http.StatusNotFound, rec.Code)
	fmt.Println(rec.Body.String())
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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

func TestAccessControlController_GetRole(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles/:roleId", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues("3")

	// when
	handleWithFilter(AccessControlController{}.GetRole, ctx)

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

func TestAccessControlController_GetRole_ID가_없는_경우(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/access-control/roles/:roleId", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("roleId")
	ctx.SetParamValues("11")

	// when
	handleWithFilter(AccessControlController{}.GetRole, ctx)

	// then
	assert.Equal(t, http.StatusNotFound, rec.Code)
	fmt.Println(rec.Body.String())
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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

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

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(AccessControlController{}.UpdateRole, ctx)

	// then
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "non changeable", resp.(map[string]interface{})["message"])
}
