package factory

import (
	"better-admin-backend-service/domain/member"
	organiztion "better-admin-backend-service/domain/organization"
	"better-admin-backend-service/dtos"
	"context"
)

type MemberAssignedAllRoleAndPermissionFactory struct {
}

func (f MemberAssignedAllRoleAndPermissionFactory) Create(ctx context.Context, member member.MemberEntity) (dtos.MemberAssignedAllRoleAndPermission, error) {
	memberAssignedAllRoleAndPermission := dtos.MemberAssignedAllRoleAndPermission{}

	filters := map[string]interface{}{}
	filters["memberId"] = member.ID
	organizationsOfMember, err := organiztion.OrganizationService{}.GetAllOrganizations(ctx, filters)
	if err != nil {
		return memberAssignedAllRoleAndPermission, nil
	}

	// 역할과 권한의 중복을 없애기 위해 MAP을 사용함.
	roleKeys := make(map[string]bool)
	assignedAllRoleNames := make([]string, 0)
	permissionKeys := make(map[string]bool)
	assignedAllPermissionNames := make([]string, 0)

	for _, role := range member.Roles {
		if _, value := roleKeys[role.Name]; !value {
			roleKeys[role.Name] = true
			assignedAllRoleNames = append(assignedAllRoleNames, role.Name)
		}

		for _, permission := range role.Permissions {
			if _, value := permissionKeys[permission.Name]; !value {
				permissionKeys[permission.Name] = true
				assignedAllPermissionNames = append(assignedAllPermissionNames, permission.Name)
			}
		}
	}

	for _, memberOrganization := range organizationsOfMember {
		for _, role := range memberOrganization.Roles {
			if _, value := roleKeys[role.Name]; !value {
				roleKeys[role.Name] = true
				assignedAllRoleNames = append(assignedAllRoleNames, role.Name)
			}

			for _, permission := range role.Permissions {
				if _, value := permissionKeys[permission.Name]; !value {
					permissionKeys[permission.Name] = true
					assignedAllPermissionNames = append(assignedAllPermissionNames, permission.Name)
				}
			}
		}
	}

	memberAssignedAllRoleAndPermission.Roles = assignedAllRoleNames
	memberAssignedAllRoleAndPermission.Permissions = assignedAllPermissionNames

	return memberAssignedAllRoleAndPermission, nil
}
