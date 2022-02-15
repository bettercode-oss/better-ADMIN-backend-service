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

func TestMenuController_CreateMenu_최상위_메뉴(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
		"name": "테스트 메뉴",
    "icon": "CiCircleOutlined",
    "link": "/test",
		"accessPermissionIds": [2, 3]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/menus", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.CreateMenu, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMenuController_CreateMenu_하위_메뉴(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	requestBody := `{
    "parentMenuId": 1,
		"name": "하위 메뉴",
    "icon": "CiCircleOutlined",
    "link": "/sub"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/menus", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.CreateMenu, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMenuController_GetMenus(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/menus", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(MenuController{}.GetMenus, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := []interface{}{
		map[string]interface{}{
			"id":       float64(1),
			"type":     "SUB_MENU",
			"name":     "최상위 메뉴 1",
			"title":    "최상위 메뉴 1",
			"icon":     "CiCircleOutlined",
			"disabled": false,
			"accessPermissions": []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "MANAGE_SYSTEM_SETTINGS",
				},
				map[string]interface{}{
					"id":   float64(2),
					"name": "MANAGE_MEMBERS",
				},
			},
			"subMenus": []interface{}{
				map[string]interface{}{
					"id":       float64(2),
					"type":     "SUB_MENU",
					"name":     "하위 메뉴 1",
					"title":    "하위 메뉴 1",
					"icon":     "CiCircleOutlined",
					"disabled": false,
					"accessPermissions": []interface{}{
						map[string]interface{}{
							"id":   float64(3),
							"name": "ACCESS_STOCK",
						},
					},
					"subMenus": []interface{}{
						map[string]interface{}{
							"id":       float64(4),
							"type":     "URL",
							"name":     "최하위 메뉴",
							"title":    "최하위 메뉴",
							"icon":     "CiCircleOutlined",
							"disabled": false,
							"link":     "/sub-sub-1",
						},
					},
				},
				map[string]interface{}{
					"id":       float64(3),
					"type":     "URL",
					"name":     "하위 메뉴 2",
					"title":    "하위 메뉴 2",
					"icon":     "CiCircleOutlined",
					"disabled": false,
					"link":     "/sub-2",
				},
			},
		},
		map[string]interface{}{
			"id":       float64(5),
			"type":     "URL",
			"title":    "최상위 메뉴2",
			"name":     "최상위 메뉴2",
			"icon":     "CiCircleOutlined",
			"disabled": false,
			"link":     "/top-2",
		},
	}

	assert.Equal(t, expected, resp.([]interface{}))
}

func TestMenuController_ChangePosition_하위로_변경(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	menuId := "2"
	requestBody := `{
		"parentMenuId": 5,
    "sameDepthMenusSequence": [2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/menus/:menuId/change-position", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("menuId")
	ctx.SetParamValues(menuId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.ChangePosition, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMenuController_ChangePosition_최상위로_변경(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	menuId := "2"
	requestBody := `{
		"sameDepthMenusSequence": [1, 2, 5]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/menu/:menuId/change-position", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("menuId")
	ctx.SetParamValues(menuId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.ChangePosition, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMenuController_DeleteMenu_최하위(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	menuId := "4"

	req := httptest.NewRequest(http.MethodDelete, "/api/menus/:menuId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("menuId")
	ctx.SetParamValues(menuId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.DeleteMenu, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMenuController_DeleteMenu_최상위(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	menuId := "1"

	req := httptest.NewRequest(http.MethodDelete, "/api/menus/:menuId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("menuId")
	ctx.SetParamValues(menuId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.DeleteMenu, ctx)

	// then
	fmt.Println(rec.Body.String())
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMenuController_UpdateMenu(t *testing.T) {
	DatabaseFixture{}.setUpDefault()

	// given
	menuId := "1"
	requestBody := `{
		"name": "바꾼 메뉴",
		"icon": "CiCircle",
		"link": "/top-2"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:menuId", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("menuId")
	ctx.SetParamValues(menuId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(MenuController{}.UpdateMenu, ctx)

	// then
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
