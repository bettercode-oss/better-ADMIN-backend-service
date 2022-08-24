package entity

import (
	memberEntity "better-admin-backend-service/domain/member/entity"
	rbacEntity "better-admin-backend-service/domain/rbac/entity"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"fmt"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type OrganizationEntity struct {
	gorm.Model
	Name                 string `gorm:"type:varchar(100);not null"`
	ParentOrganizationID *uint
	ParentOrganization   *OrganizationEntity
	Path                 string                      `gorm:"-"`
	Roles                []rbacEntity.RoleEntity     `gorm:"many2many:organization_roles;"`
	Members              []memberEntity.MemberEntity `gorm:"many2many:organization_members;"`
	CreatedBy            uint
	UpdatedBy            uint
}

func (OrganizationEntity) TableName() string {
	return "organizations"
}

func (o *OrganizationEntity) ChangePosition(ctx context.Context, parentOrganizationId *uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	o.ParentOrganizationID = parentOrganizationId
	o.UpdatedBy = userClaim.Id

	return nil
}

func (o *OrganizationEntity) GeneratePath(entities []OrganizationEntity) {
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

func (o OrganizationEntity) FindChildEntities(entities []OrganizationEntity) ([]OrganizationEntity, error) {
	childEntities := make([]OrganizationEntity, 0)

	for i := 0; i < len(entities); i++ {
		entities[i].GeneratePath(entities)
		if strings.Contains(entities[i].Path, strconv.FormatUint(uint64(o.ID), 10)) {
			childEntities = append(childEntities, entities[i])
		}
	}

	return childEntities, nil
}

func (o *OrganizationEntity) AssignRole(ctx context.Context, roleEntities []rbacEntity.RoleEntity) error {
	// 기존 역할을 덮어쓰기
	o.Roles = roleEntities

	return nil
}

func (o *OrganizationEntity) AssignMember(ctx context.Context, memberEntities []memberEntity.MemberEntity) error {
	o.Members = memberEntities

	return nil
}

func (o *OrganizationEntity) ChangeName(ctx context.Context, name string) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	o.Name = name
	o.UpdatedBy = userClaim.Id
	return nil
}

func (o OrganizationEntity) ExistMember(memberId uint) bool {
	for _, member := range o.Members {
		if member.ID == memberId {
			return true
		}
	}

	return false
}

func NewOrganizationEntity(ctx context.Context, information dtos.OrganizationInformation) (OrganizationEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return OrganizationEntity{}, err
	}

	return OrganizationEntity{
		Name:                 information.Name,
		ParentOrganizationID: information.ParentOrganizationId,
		CreatedBy:            userClaim.Id,
		UpdatedBy:            userClaim.Id,
	}, nil
}
