package dtos

import (
	"github.com/labstack/echo"
)

const (
	TypeMenuUrl     = "URL"
	TypeMenuSubMenu = "SUB_MENU"
)

type MenuInformation struct {
	Id                  uint                   `json:"id"`
	Type                string                 `json:"type"`
	Title               string                 `json:"title"`
	Name                string                 `json:"name" validate:"required"`
	Icon                string                 `json:"icon"`
	Link                *string                `json:"link,omitempty"`
	Disabled            bool                   `json:"disabled"`
	ParentMenuId        *uint                  `json:"parentMenuId,omitempty"`
	SubMenus            []MenuInformation      `json:"subMenus,omitempty"`
	AccessPermissionIds []uint                 `json:"accessPermissionIds,omitempty"`
	AccessPermissions   []MenuAccessPermission `json:"accessPermissions,omitempty"`
}

func (m MenuInformation) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}

func (m *MenuInformation) SetUpMenuType() {
	if m.Link != nil && len(*m.Link) > 0 {
		m.Type = TypeMenuUrl
	} else {
		m.Type = TypeMenuSubMenu
	}
}

type MenuPosition struct {
	ParentMenuId           *uint  `json:"parentMenuId"`
	SameDepthMenusSequence []uint `json:"sameDepthMenusSequence"`
}

type MenuAccessPermission struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
