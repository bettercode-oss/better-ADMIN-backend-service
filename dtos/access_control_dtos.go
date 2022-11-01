package dtos

import (
	"github.com/labstack/echo"
	"time"
)

type PermissionInformation struct {
	Id          uint   `json:"id"`
	Type        string `json:"type"`
	TypeName    string `json:"typeName"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (p PermissionInformation) Validate(ctx echo.Context) error {
	return ctx.Validate(p)
}

type PermissionDetails struct {
	Id          uint      `json:"id"`
	Type        string    `json:"type"`
	TypeName    string    `json:"typeName"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RoleInformation struct {
	Name                 string `json:"name" validate:"required"`
	Description          string `json:"description"`
	AllowedPermissionIds []uint `json:"allowedPermissionIds" validate:"required"`
}

func (r RoleInformation) Validate(ctx echo.Context) error {
	return ctx.Validate(r)
}

type RoleSummary struct {
	Id                uint                `json:"id"`
	Type              string              `json:"type"`
	TypeName          string              `json:"typeName"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	AllowedPermission []AllowedPermission `json:"permissions"`
}

type AllowedPermission struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type RoleDetails struct {
	Id                 uint                `json:"id"`
	Type               string              `json:"type"`
	TypeName           string              `json:"typeName"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	CreatedAt          time.Time           `json:"createdAt"`
	AllowedPermissions []AllowedPermission `json:"permissions"`
}
