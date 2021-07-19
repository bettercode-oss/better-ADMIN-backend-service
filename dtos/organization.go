package dtos

import (
	"github.com/labstack/echo"
)

type OrganizationInformation struct {
	Id                   uint                      `json:"id"`
	Name                 string                    `json:"name" validate:"required"`
	ParentOrganizationId *uint                     `json:"parentOrganizationId,omitempty"`
	SubOrganizations     []OrganizationInformation `json:"subOrganizations,omitempty"`
	OrganizationRoles    []OrganizationRole        `json:"roles,omitempty"`
	OrganizationMembers  []OrganizationMember      `json:"members,omitempty"`
}

func (o OrganizationInformation) Validate(ctx echo.Context) error {
	return ctx.Validate(o)
}

type OrganizationAssignRole struct {
	RoleIds []uint `json:"roleIds" validate:"required"`
}

func (o OrganizationAssignRole) Validate(ctx echo.Context) error {
	return ctx.Validate(o)
}

type OrganizationRole struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type OrganizationMember struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type OrganizationAssignMember struct {
	MemberIds []uint `json:"memberIds" validate:"required"`
}

func (o OrganizationAssignMember) Validate(ctx echo.Context) error {
	return ctx.Validate(o)
}
