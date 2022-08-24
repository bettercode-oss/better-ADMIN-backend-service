package entity

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"gorm.io/gorm"
)

const (
	PreDefineTypeKey   = "pre-define"
	PreDefineTypeName  = "사전정의"
	UserDefineTypeKey  = "user-define"
	UserDefineTypeName = "사용자정의"
)

type PermissionEntity struct {
	gorm.Model
	Type        string `gorm:"type:varchar(50);not null"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(1000)"`
	CreatedBy   uint
	UpdatedBy   uint
}

func (PermissionEntity) TableName() string {
	return "permissions"
}

func (p PermissionEntity) GetTypeName() string {
	if p.Type == PreDefineTypeKey {
		return PreDefineTypeName
	}

	if p.Type == UserDefineTypeKey {
		return UserDefineTypeName
	}

	return ""
}

func (p *PermissionEntity) Update(ctx context.Context, information dtos.PermissionInformation) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	if p.Type == PreDefineTypeKey {
		return domain.ErrNonChangeable
	}

	p.Name = information.Name
	p.Description = information.Description
	p.UpdatedBy = userClaim.Id

	return nil
}

func (p PermissionEntity) Deletable() error {
	if p.Type == PreDefineTypeKey {
		return domain.ErrNonChangeable
	}

	return nil
}

func NewPermissionEntity(ctx context.Context, information dtos.PermissionInformation) (PermissionEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return PermissionEntity{}, err
	}

	return PermissionEntity{
		Type:        UserDefineTypeKey,
		Name:        information.Name,
		Description: information.Description,
		CreatedBy:   userClaim.Id,
		UpdatedBy:   userClaim.Id,
	}, nil
}

type RoleEntity struct {
	gorm.Model
	Type        string `gorm:"type:varchar(50);not null"`
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(1000)"`
	CreatedBy   uint
	UpdatedBy   uint
	Permissions []PermissionEntity `gorm:"many2many:role_permissions;"`
}

func (RoleEntity) TableName() string {
	return "roles"
}

func (r RoleEntity) GetTypeName() string {
	if r.Type == PreDefineTypeKey {
		return PreDefineTypeName
	}

	if r.Type == UserDefineTypeKey {
		return UserDefineTypeName
	}

	return ""
}

func (r RoleEntity) Deletable() error {
	if r.Type == PreDefineTypeKey {
		return domain.ErrNonChangeable
	}

	return nil
}

func (r *RoleEntity) Update(ctx context.Context, information dtos.RoleInformation, permissionEntities []PermissionEntity) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	if r.Type == PreDefineTypeKey {
		return domain.ErrNonChangeable
	}

	r.Name = information.Name
	r.Description = information.Description
	r.UpdatedBy = userClaim.Id
	r.Permissions = permissionEntities
	return nil
}
