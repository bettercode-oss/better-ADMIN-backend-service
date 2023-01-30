package services

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	memberDomain "better-admin-backend-service/member/domain"
	"better-admin-backend-service/organization/domain"
	"better-admin-backend-service/organization/repository"
	"context"
	"github.com/wesovilabs/koazee"
	"strings"
)

type OrganizationService struct {
	rbacService            *RoleBasedAccessControlService
	organizationRepository *repository.OrganizationRepository
	memberService          *MemberService
}

func NewOrganizationService(
	rbacService *RoleBasedAccessControlService,
	organizationRepository *repository.OrganizationRepository,
	memberService *MemberService) *OrganizationService {
	return &OrganizationService{
		rbacService:            rbacService,
		organizationRepository: organizationRepository,
		memberService:          memberService,
	}
}

func (s OrganizationService) CreateOrganization(ctx context.Context, information dtos.OrganizationInformation) error {
	organizationEntity, err := domain.NewOrganizationEntity(ctx, information)
	if err != nil {
		return err
	}
	return s.organizationRepository.Create(ctx, organizationEntity)
}

func (s OrganizationService) GetAllOrganizations(ctx context.Context, filters map[string]interface{}) ([]domain.OrganizationEntity, error) {
	entities, err := s.organizationRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].GeneratePath(entities)
	}

	entitiesSortedByPath := koazee.StreamOf(entities).Sort(func(a, b domain.OrganizationEntity) int {
		return strings.Compare(a.Path, b.Path)
	}).Out().Val().([]domain.OrganizationEntity)

	return entitiesSortedByPath, nil
}

func (s OrganizationService) ChangePosition(ctx context.Context, organizationId uint, parentOrganizationId *uint) error {
	organizationEntity, err := s.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.ChangePosition(ctx, parentOrganizationId)
	if err != nil {
		return err
	}

	return s.organizationRepository.Save(ctx, &organizationEntity)
}

func (s OrganizationService) DeleteOrganization(ctx context.Context, organizationId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	organizationEntity, err := s.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	entities, err := s.organizationRepository.FindAll(ctx, nil)
	if err != nil {
		return err
	}

	childEntities, err := organizationEntity.FindChildEntities(entities)
	if err != nil {
		return err
	}

	for _, childEntity := range childEntities {
		childEntity.UpdatedBy = userClaim.Id
		if err := s.organizationRepository.Delete(ctx, childEntity); err != nil {
			return err
		}
	}

	organizationEntity.UpdatedBy = userClaim.Id
	return s.organizationRepository.Delete(ctx, organizationEntity)
}

func (s OrganizationService) AssignRoles(ctx context.Context, organizationId uint, assignRole dtos.OrganizationAssignRole) error {
	organizationEntity, err := s.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["roleIds"] = assignRole.RoleIds

	findRoleEntities, _, err := s.rbacService.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = organizationEntity.AssignRole(ctx, findRoleEntities)
	if err != nil {
		return err
	}

	return s.organizationRepository.Save(ctx, &organizationEntity)
}

func (s OrganizationService) AssignMembers(ctx context.Context, organizationId uint, assignMember dtos.OrganizationAssignMember) error {
	organizationEntity, err := s.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{}
	filters["memberIds"] = assignMember.MemberIds

	findMemberEntities, _, err := s.memberService.GetMembers(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	err = organizationEntity.AssignMember(ctx, findMemberEntities)
	if err != nil {
		return err
	}

	return s.organizationRepository.Save(ctx, &organizationEntity)
}

func (s OrganizationService) ChangeOrganizationName(ctx context.Context, organizationId uint, organizationName string) error {
	organizationEntity, err := s.organizationRepository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.ChangeName(ctx, organizationName)
	if err != nil {
		return err
	}

	return s.organizationRepository.Save(ctx, &organizationEntity)
}

func (s OrganizationService) GetMemberAssignedAllRoleAndPermission(ctx context.Context, member memberDomain.MemberEntity) (dtos.MemberAssignedAllRoleAndPermission, error) {
	memberAssignedAllRoleAndPermission := dtos.MemberAssignedAllRoleAndPermission{}

	filters := map[string]interface{}{}
	filters["memberId"] = member.ID
	organizationsOfMember, err := s.GetAllOrganizations(ctx, filters)
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

func (s OrganizationService) GetOrganization(ctx context.Context, organizationId uint) (domain.OrganizationEntity, error) {
	return s.organizationRepository.FindById(ctx, organizationId)
}
