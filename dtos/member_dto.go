package dtos

import (
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
	LastAccessAt        *time.Time           `json:"lastAccessAt"`
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
	RoleIds []uint `json:"roleIds" binding:"required"`
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
	SignId   string `json:"signId" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
