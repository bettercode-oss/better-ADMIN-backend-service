package services

import (
	memberEntity "better-admin-backend-service/domain/member/entity"
	"better-admin-backend-service/domain/organization/entity"
	"better-admin-backend-service/domain/organization/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/wesovilabs/koazee"
	"strings"
)

type OrganizationService struct {
}

func (OrganizationService) CreateOrganization(ctx context.Context, information dtos.OrganizationInformation) error {
	organizationEntity, err := entity.NewOrganizationEntity(ctx, information)
	if err != nil {
		return err
	}
	return repository.OrganizationRepository{}.Create(ctx, organizationEntity)
}

func (OrganizationService) GetAllOrganizations(ctx context.Context, filters map[string]interface{}) ([]entity.OrganizationEntity, error) {
	entities, err := repository.OrganizationRepository{}.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].GeneratePath(entities)
	}

	entitiesSortedByPath := koazee.StreamOf(entities).Sort(func(a, b entity.OrganizationEntity) int {
		return strings.Compare(a.Path, b.Path)
	}).Out().Val().([]entity.OrganizationEntity)

	return entitiesSortedByPath, nil
}

func (OrganizationService) ChangePosition(ctx context.Context, organizationId uint, parentOrganizationId *uint) error {
	repository := repository.OrganizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.ChangePosition(ctx, parentOrganizationId)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (OrganizationService) DeleteOrganization(ctx context.Context, organizationId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := repository.OrganizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	entities, err := repository.FindAll(ctx, nil)
	if err != nil {
		return err
	}

	childEntities, err := organizationEntity.FindChildEntities(entities)
	if err != nil {
		return err
	}

	for _, childEntity := range childEntities {
		childEntity.UpdatedBy = userClaim.Id
		if err := repository.Delete(ctx, childEntity); err != nil {
			return err
		}
	}

	organizationEntity.UpdatedBy = userClaim.Id
	return repository.Delete(ctx, organizationEntity)
}

func (OrganizationService) AssignRoles(ctx context.Context, organizationId uint, assignRole dtos.OrganizationAssignRole) error {
	repository := repository.OrganizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["roleIds"] = assignRole.RoleIds

	findRoleEntities, _, err := RoleBasedAccessControlService{}.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = organizationEntity.AssignRole(ctx, findRoleEntities)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (OrganizationService) AssignMembers(ctx context.Context, organizationId uint, assignMember dtos.OrganizationAssignMember) error {
	repository := repository.OrganizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["memberIds"] = assignMember.MemberIds

	findMemberEntities, _, err := MemberService{}.GetMembers(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = organizationEntity.AssignMember(ctx, findMemberEntities)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (OrganizationService) ChangeOrganizationName(ctx context.Context, organizationId uint, organizationName string) error {
	repository := repository.OrganizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.ChangeName(ctx, organizationName)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (service OrganizationService) GetMemberAssignedAllRoleAndPermission(ctx context.Context, member memberEntity.MemberEntity) (dtos.MemberAssignedAllRoleAndPermission, error) {
	memberAssignedAllRoleAndPermission := dtos.MemberAssignedAllRoleAndPermission{}

	filters := map[string]interface{}{}
	filters["memberId"] = member.ID
	organizationsOfMember, err := service.GetAllOrganizations(ctx, filters)
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

func (OrganizationService) GetOrganization(ctx context.Context, organizationId uint) (entity.OrganizationEntity, error) {
	return repository.OrganizationRepository{}.FindById(ctx, organizationId)
}
