package organization

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/wesovilabs/koazee"
	"strings"
)

type OrganizationService struct {
}

func (OrganizationService) CreateOrganization(ctx context.Context, information dtos.OrganizationInformation) error {
	organizationEntity, err := NewOrganizationEntity(ctx, information)
	if err != nil {
		return err
	}
	return organizationRepository{}.Create(ctx, organizationEntity)
}

func (OrganizationService) GetAllOrganizations(ctx context.Context, filters map[string]interface{}) ([]OrganizationEntity, error) {
	entities, err := organizationRepository{}.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].generatePath(entities)
	}

	entitiesSortedByPath := koazee.StreamOf(entities).Sort(func(a, b OrganizationEntity) int {
		return strings.Compare(a.Path, b.Path)
	}).Out().Val().([]OrganizationEntity)

	return entitiesSortedByPath, nil
}

func (OrganizationService) ChangePosition(ctx context.Context, organizationId uint, parentOrganizationId *uint) error {
	repository := organizationRepository{}
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

	repository := organizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	childEntities, err := organizationEntity.FindChildEntities(ctx)
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
	repository := organizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.AssignRole(ctx, assignRole)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (OrganizationService) AssignMembers(ctx context.Context, organizationId uint, assignMember dtos.OrganizationAssignMember) error {
	repository := organizationRepository{}
	organizationEntity, err := repository.FindById(ctx, organizationId)
	if err != nil {
		return err
	}

	err = organizationEntity.AssignMember(ctx, assignMember)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &organizationEntity)
}

func (OrganizationService) ChangeOrganizationName(ctx context.Context, organizationId uint, organizationName string) error {
	repository := organizationRepository{}
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
