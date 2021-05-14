package member

import (
	"better-admin-backend-service/domain/rbac"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestMemberEntity_GetPermissionNames(t *testing.T) {
	// given
	entity := MemberEntity{
		Model:  gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Type:   "site",
		SignId: "ymyoo",
		Name:   "유영모",
		Roles: []rbac.RoleEntity{
			{
				Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Type:  "user-define",
				Name:  "테스터",
				Permissions: []rbac.PermissionEntity{
					{
						Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Type:  "user-define",
						Name:  "권한1",
					},
					{
						Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Type:  "user-define",
						Name:  "권한2",
					},
				},
			},
			{
				Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Type:  "user-define",
				Name:  "테스터2",
				Permissions: []rbac.PermissionEntity{
					{
						Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Type:  "user-define",
						Name:  "권한2",
					},
					{
						Model: gorm.Model{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Type:  "user-define",
						Name:  "권한3",
					},
				},
			},
		},
	}

	// when
	permissionNames := entity.GetPermissionNames()

	// then
	assert.Equal(t, []string{"권한1", "권한2", "권한3"}, permissionNames)
}
