package dtos

import (
	"time"
)

type OrganizationInformation struct {
	Id                   uint                      `json:"id"`
	Name                 string                    `json:"name" binding:"required"`
	ParentOrganizationId *uint                     `json:"parentOrganizationId,omitempty"`
	SubOrganizations     []OrganizationInformation `json:"subOrganizations,omitempty"`
	OrganizationRoles    []OrganizationRole        `json:"roles,omitempty"`
	OrganizationMembers  []OrganizationMember      `json:"members,omitempty"`
}

type OrganizationAssignRole struct {
	RoleIds []uint `json:"roleIds" binding:"required"`
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
	MemberIds []uint `json:"memberIds" binding:"required"`
}

type OrganizationDetails struct {
	Id        uint                 `json:"id"`
	Name      string               `json:"name"`
	CreatedAt time.Time            `json:"createdAt"`
	Roles     []OrganizationRole   `json:"roles,omitempty"`
	Members   []OrganizationMember `json:"members,omitempty"`
}
