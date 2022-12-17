package dtos

import (
	"time"
)

type PermissionInformation struct {
	Id          uint   `json:"id"`
	Type        string `json:"type"`
	TypeName    string `json:"typeName"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
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
	Name                 string `json:"name" binding:"required"`
	Description          string `json:"description"`
	AllowedPermissionIds []uint `json:"allowedPermissionIds" binding:"required"`
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
