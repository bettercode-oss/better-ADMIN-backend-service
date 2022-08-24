package factory

import (
	"better-admin-backend-service/domain/organization/entity"
	"better-admin-backend-service/dtos"
)

func NewOrganizationInformationFromEntity(entity entity.OrganizationEntity) dtos.OrganizationInformation {
	organizationInformation := dtos.OrganizationInformation{
		Id:   entity.ID,
		Name: entity.Name,
	}

	if entity.Roles != nil && len(entity.Roles) > 0 {
		roles := make([]dtos.OrganizationRole, 0)
		for _, role := range entity.Roles {
			roles = append(roles, dtos.OrganizationRole{
				Id:   role.ID,
				Name: role.Name,
			})
		}
		organizationInformation.OrganizationRoles = roles
	}

	if entity.Members != nil && len(entity.Members) > 0 {
		members := make([]dtos.OrganizationMember, 0)
		for _, member := range entity.Members {
			members = append(members, dtos.OrganizationMember{
				Id:   member.ID,
				Name: member.Name,
			})
		}
		organizationInformation.OrganizationMembers = members
	}

	return organizationInformation
}
