package dtos

import (
	"github.com/labstack/echo"
)

type MemberInformation struct {
	Id          uint         `json:"id"`
	Type        string       `json:"type"`
	TypeName    string       `json:"typeName"`
	Name        string       `json:"name"`
	MemberRoles []MemberRole `json:"roles"`
}

type MemberRole struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type MemberAssignRole struct {
	RoleIds []uint `json:"roleIds" validate:"required"`
}

func (r MemberAssignRole) Validate(ctx echo.Context) error {
	return ctx.Validate(r)
}

type CurrentMember struct {
	Id          uint     `json:"id"`
	Type        string   `json:"type"`
	TypeName    string   `json:"typeName"`
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
