package member

import (
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/dtos"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	TypeMemberSite       = "site"
	TypeMemberSiteName   = "사이트"
	TypeMemberDooray     = "dooray"
	TypeMemberDoorayName = "두레이"
)

type MemberEntity struct {
	gorm.Model
	Type           string
	SignId         string
	Name           string
	Password       string
	DoorayId       string
	DoorayUserCode string
	Roles          []rbac.RoleEntity `gorm:"many2many:member_roles;"`
}

func (MemberEntity) TableName() string {
	return "members"
}

func (m MemberEntity) ValidatePassword(password string) error {
	if m.comparePasswords(m.Password, password) == false {
		return errors.New("InvalidPassword")
	}

	return nil
}

func (m MemberEntity) hashAndSalt(pwd string) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func (m MemberEntity) comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		return false
	}

	return true
}

func (m MemberEntity) GetTypeName() string {
	if m.Type == TypeMemberSite {
		return TypeMemberSiteName
	}

	if m.Type == TypeMemberDooray {
		return TypeMemberDoorayName
	}

	return ""
}

func (m *MemberEntity) AssignRole(ctx context.Context, role dtos.MemberAssignRole) error {
	// 기존 역할을 덮어쓰기
	filters := map[string]interface{}{}
	filters["roleIds"] = role.RoleIds

	findRoleEntities, _, err := rbac.RoleBasedAccessControlService{}.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	m.Roles = findRoleEntities

	return nil
}

func (m MemberEntity) GetRoleNames() []string {
	var rolesNames = make([]string, 0)
	if m.Roles == nil {
		return rolesNames
	}

	for _, role := range m.Roles {
		rolesNames = append(rolesNames, role.Name)
	}

	return rolesNames
}

func (m MemberEntity) GetPermissionNames() []string {
	// 역할에 할당된 권한을 반환한다.
	// 권한이 중복이 일어날 수 있기 때문에 중복을 없애고 반환한다.
	keys := make(map[string]bool)
	permissionNames := make([]string, 0)
	if m.Roles == nil {
		return permissionNames
	}

	for _, role := range m.Roles {
		if role.Permissions == nil {
			continue
		}

		for _, permission := range role.Permissions {
			if _, value := keys[permission.Name]; !value {
				keys[permission.Name] = true
				permissionNames = append(permissionNames, permission.Name)
			}
		}
	}

	return permissionNames
}
