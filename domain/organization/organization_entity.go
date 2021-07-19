package organization

import (
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/dtos"
	"context"
	"fmt"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type OrganizationEntity struct {
	gorm.Model
	Name                 string
	ParentOrganizationID *uint
	ParentOrganization   *OrganizationEntity
	Path                 string                `gorm:"-"`
	Roles                []rbac.RoleEntity     `gorm:"many2many:organization_roles;"`
	Members              []member.MemberEntity `gorm:"many2many:organization_members;"`
}

func (OrganizationEntity) TableName() string {
	return "organizations"
}

func (o *OrganizationEntity) ChangePosition(parentOrganizationId *uint) {
	o.ParentOrganizationID = parentOrganizationId
}

func (o *OrganizationEntity) generatePath(entities []OrganizationEntity) {
	fullPath := o.getPath(o.ID, entities, "")
	o.Path = strings.Join(koazee.StreamOf(strings.Split(fullPath, "-")).Reverse().Out().Val().([]string), "-")
}

func (o OrganizationEntity) getPath(targetId uint, organizations []OrganizationEntity, path string) string {
	for _, en := range organizations {
		if en.ID == targetId {
			if en.ParentOrganizationID == nil {
				return path
			}
			if path == "" {
				path = fmt.Sprintf("%v", *en.ParentOrganizationID)
			} else {
				path += fmt.Sprintf("-%v", *en.ParentOrganizationID)
			}

			return o.getPath(*en.ParentOrganizationID, organizations, path)
		}
	}
	return ""
}

func (o OrganizationEntity) FindChildEntities(ctx context.Context) ([]OrganizationEntity, error) {
	childEntities := make([]OrganizationEntity, 0)

	entities, err := organizationRepository{}.FindAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].generatePath(entities)
		if strings.Contains(entities[i].Path, strconv.FormatUint(uint64(o.ID), 10)) {
			childEntities = append(childEntities, entities[i])
		}
	}

	return childEntities, nil
}

func (o *OrganizationEntity) AssignRole(ctx context.Context, role dtos.OrganizationAssignRole) error {
	// 기존 역할을 덮어쓰기
	filters := map[string]interface{}{}
	filters["roleIds"] = role.RoleIds

	findRoleEntities, _, err := rbac.RoleBasedAccessControlService{}.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	o.Roles = findRoleEntities

	return nil
}

func (o *OrganizationEntity) AssignMember(ctx context.Context, assignMember dtos.OrganizationAssignMember) error {
	filters := map[string]interface{}{}
	filters["memberIds"] = assignMember.MemberIds

	findMemberEntities, _, err := member.MemberService{}.GetMembers(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	o.Members = findMemberEntities

	return nil
}

func (o *OrganizationEntity) ChangeName(name string) {
	o.Name = name
}

func (o OrganizationEntity) ExistMember(memberId uint) bool {
	for _, member := range o.Members {
		if member.ID == memberId {
			return true
		}
	}

	return false
}

func NewOrganizationEntity(information dtos.OrganizationInformation) OrganizationEntity {
	return OrganizationEntity{
		Name:                 information.Name,
		ParentOrganizationID: information.ParentOrganizationId,
	}
}
