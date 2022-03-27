package controllers

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/menu"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/factory"
	"better-admin-backend-service/middlewares"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type MenuController struct {
}

func (controller MenuController) Init(g *echo.Group) {
	g.POST("", controller.CreateMenu, middlewares.CheckPermission([]string{domain.PermissionManageMenus}))
	g.GET("", controller.GetMenus, middlewares.CheckPermission([]string{"*"}))
	g.PUT("/:menuId/change-position", controller.ChangePosition, middlewares.CheckPermission([]string{domain.PermissionManageMenus}))
	g.PUT("/:menuId", controller.UpdateMenu, middlewares.CheckPermission([]string{domain.PermissionManageMenus}))
	g.DELETE("/:menuId", controller.DeleteMenu, middlewares.CheckPermission([]string{domain.PermissionManageMenus}))
}

func (MenuController) CreateMenu(ctx echo.Context) error {
	var menuInformation dtos.MenuInformation
	if err := ctx.Bind(&menuInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := menuInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := menu.MenuService{}.CreateMenu(ctx.Request().Context(), menuInformation)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusCreated)
}

func (MenuController) GetMenus(ctx echo.Context) error {
	allOfMenus, err := menu.MenuService{}.GetAllMenus(ctx.Request().Context())
	if err != nil {
		return err
	}

	menus := make([]dtos.MenuInformation, 0)
	for _, entity := range allOfMenus {
		if entity.ParentMenuId == nil {
			menuInformation := factory.NewMenuInformationFromEntity(entity)
			menus = append(menus, menuInformation)
			continue
		}

		parentMenuInformation := findParentMenuInformation(&menus, *entity.ParentMenuId)
		if parentMenuInformation == nil {
			return errors.New("not found parentMenuInformation")
		}

		if parentMenuInformation.SubMenus == nil {
			parentMenuInformation.SubMenus = make([]dtos.MenuInformation, 0)
		}

		menusAccessPermission := make([]dtos.MenuAccessPermission, 0)
		for _, permission := range entity.Permissions {
			menusAccessPermission = append(menusAccessPermission, dtos.MenuAccessPermission{
				Id:   permission.ID,
				Name: permission.Name,
			})
		}
		menuInformation := factory.NewMenuInformationFromEntity(entity)
		parentMenuInformation.SubMenus = append(parentMenuInformation.SubMenus, menuInformation)
	}

	return ctx.JSON(http.StatusOK, menus)
}

func findParentMenuInformation(menus *[]dtos.MenuInformation, parentId uint) *dtos.MenuInformation {
	for i := 0; i < len(*menus); i++ {
		if (*menus)[i].Id == parentId {
			return &(*menus)[i]
		}
		if (*menus)[i].SubMenus != nil {
			find := findParentMenuInformation(&(*menus)[i].SubMenus, parentId)
			if find != nil {
				return find
			}
		}
	}
	return nil
}

func (MenuController) ChangePosition(ctx echo.Context) error {
	menuId, err := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var menuPosition dtos.MenuPosition
	if err := ctx.Bind(&menuPosition); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = menu.MenuService{}.ChangePosition(ctx.Request().Context(), uint(menuId), menuPosition)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (MenuController) DeleteMenu(ctx echo.Context) error {
	menuId, err := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = menu.MenuService{}.DeleteMenu(ctx.Request().Context(), uint(menuId))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (MenuController) UpdateMenu(ctx echo.Context) error {
	menuId, err := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var menuInformation dtos.MenuInformation
	if err := ctx.Bind(&menuInformation); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := menuInformation.Validate(ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err = menu.MenuService{}.UpdateMenu(ctx.Request().Context(), uint(menuId), menuInformation)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
