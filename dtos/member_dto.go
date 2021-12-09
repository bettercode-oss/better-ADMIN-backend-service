package dtos

import (
	"github.com/labstack/echo"
	"time"
)

type MemberInformation struct {
	Id                  uint                 `json:"id"`
	SignId              string               `json:"signId"`
	Type                string               `json:"type"`
	TypeName            string               `json:"typeName"`
	CandidateId         string               `json:"candidateId"`
	Name                string               `json:"name"`
	MemberRoles         []MemberRole         `json:"roles"`
	MemberOrganizations []MemberOrganization `json:"organizations"`
	CreatedAt           time.Time            `json:"createdAt"`
}

type MemberRole struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type MemberOrganization struct {
	Id    uint                     `json:"id"`
	Name  string                   `json:"name"`
	Roles []MemberOrganizationRole `json:"roles"`
}

type MemberOrganizationRole struct {
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
	Picture     string   `json:"picture"`
}

type MemberAssignedAllRoleAndPermission struct {
	Roles       []string
	Permissions []string
}

type MemberSignUp struct {
	SignId   string `json:"signId" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (m MemberSignUp) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}
