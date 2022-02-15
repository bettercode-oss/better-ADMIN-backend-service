package menu

import (
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"fmt"
	"github.com/wesovilabs/koazee"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type MenuEntity struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(100);not null"`
	Icon         string  `gorm:"type:varchar(50)"`
	Link         *string `gorm:"type:varchar(1000)"`
	Disabled     bool
	ParentMenuId *uint
	Sequence     uint
	Permissions  []rbac.PermissionEntity `gorm:"many2many:menu_permissions;"`
	CreatedBy    uint
	UpdatedBy    uint
	Path         string `gorm:"-"`
}

func (MenuEntity) TableName() string {
	return "menus"
}

func (m *MenuEntity) generatePath(entities []MenuEntity) {
	fullPath := m.getPath(m.ID, entities, "")
	m.Path = strings.Join(koazee.StreamOf(strings.Split(fullPath, "-")).Reverse().Out().Val().([]string), "-")
}

func (o MenuEntity) getPath(targetId uint, menus []MenuEntity, path string) string {
	for _, en := range menus {
		if en.ID == targetId {
			if en.ParentMenuId == nil {
				return path
			}
			if path == "" {
				path = fmt.Sprintf("%v", *en.ParentMenuId)
			} else {
				path += fmt.Sprintf("-%v", *en.ParentMenuId)
			}

			return o.getPath(*en.ParentMenuId, menus, path)
		}
	}
	return ""
}

func (m *MenuEntity) ChangePosition(ctx context.Context, position dtos.MenuPosition) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	m.ParentMenuId = position.ParentMenuId
	m.UpdatedBy = userClaim.Id

	return nil
}

func (m MenuEntity) FindChildEntities(ctx context.Context) ([]MenuEntity, error) {
	childEntities := make([]MenuEntity, 0)

	entities, err := menuRepository{}.FindAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].generatePath(entities)
		if strings.Contains(entities[i].Path, strconv.FormatUint(uint64(m.ID), 10)) {
			childEntities = append(childEntities, entities[i])
		}
	}

	return childEntities, nil
}

func (m *MenuEntity) Update(ctx context.Context, information dtos.MenuInformation) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	m.Name = information.Name
	m.Icon = information.Icon
	m.Disabled = information.Disabled
	m.Link = information.Link
	m.UpdatedBy = userClaim.Id

	filters := map[string]interface{}{}
	filters["permissionIds"] = information.AccessPermissionIds

	permissionEntities, _, err := rbac.RoleBasedAccessControlService{}.GetPermissions(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	m.Permissions = permissionEntities
	return nil
}

func NewMenuEntity(ctx context.Context, information dtos.MenuInformation) (MenuEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return MenuEntity{}, err
	}

	lastSequence, err := getLastMenuSequence(ctx, information)
	if err != nil {
		return MenuEntity{}, err
	}

	entity := MenuEntity{
		Name:         information.Name,
		Icon:         information.Icon,
		Disabled:     information.Disabled,
		Link:         information.Link,
		ParentMenuId: information.ParentMenuId,
		CreatedBy:    userClaim.Id,
		UpdatedBy:    userClaim.Id,
		Sequence:     lastSequence,
	}

	filters := map[string]interface{}{}
	filters["permissionIds"] = information.AccessPermissionIds

	permissionEntities, _, err := rbac.RoleBasedAccessControlService{}.GetPermissions(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return entity, err
	}

	entity.Permissions = permissionEntities
	return entity, nil
}

func getLastMenuSequence(ctx context.Context, information dtos.MenuInformation) (uint, error) {
	filters := map[string]interface{}{}
	filters["parentMenuId"] = information.ParentMenuId

	menuEntities, err := menuRepository{}.FindAll(ctx, filters)
	if err != nil {
		return 0, err
	}

	count := len(menuEntities)
	if menuEntities == nil || count == 0 {
		return 0, nil
	}

	return menuEntities[count-1].Sequence + 1, nil
}
