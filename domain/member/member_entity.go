package member

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
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
	TypeMemberGoogle     = "google"
	TypeMemberGoogleName = "구글"
	StatusMemberApplied  = "applied"
	StatusMemberApproved = "approved"
)

type MemberEntity struct {
	gorm.Model
	Type           string `gorm:"type:varchar(20);not null"`
	SignId         string `gorm:"type:varchar(50)"`
	Name           string `gorm:"type:varchar(50)"`
	Password       string `gorm:"type:varchar(100)"`
	Status         string `gorm:"type:varchar(20);not null"`
	DoorayId       string `gorm:"type:varchar(50)"`
	DoorayUserCode string `gorm:"type:varchar(50)"`
	GoogleId       string `gorm:"type:varchar(50)"`
	GoogleMail     string `gorm:"type:varchar(50)"`
	Picture        string `gorm:"type:varchar(1000)"`
	UpdatedBy      uint
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

	if m.Type == TypeMemberGoogle {
		return TypeMemberGoogleName
	}

	return ""
}

func (m *MemberEntity) AssignRole(ctx context.Context, role dtos.MemberAssignRole) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	// 기존 역할을 덮어쓰기
	filters := map[string]interface{}{}
	filters["roleIds"] = role.RoleIds

	findRoleEntities, _, err := rbac.RoleBasedAccessControlService{}.GetRoles(ctx, filters, dtos.Pageable{Page: 0})
	if err != nil {
		return err
	}

	m.Roles = findRoleEntities
	m.UpdatedBy = userClaim.Id

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

func (m *MemberEntity) Approve(ctx context.Context) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	if m.Status == StatusMemberApproved {
		return domain.ErrAlreadyApproved
	}
	m.Status = StatusMemberApproved
	m.UpdatedBy = userClaim.Id
	return nil
}

func (m MemberEntity) IsApproved() bool {
	if m.Status == StatusMemberApproved {
		return true
	}
	return false
}

func (m MemberEntity) GetCandidateId() string {
	if m.Type == TypeMemberSite {
		return m.SignId
	} else if m.Type == TypeMemberDooray {
		return m.DoorayUserCode
	} else if m.Type == TypeMemberGoogle {
		return m.GoogleMail
	} else {
		return ""
	}
}

func NewMemberEntityFromSignUp(signUp dtos.MemberSignUp) (MemberEntity, error) {
	hashedPassword, err := MemberEntity{}.hashAndSalt(signUp.Password)
	if err != nil {
		return MemberEntity{}, err
	}

	return MemberEntity{
		Type:     TypeMemberSite,
		SignId:   signUp.SignId,
		Name:     signUp.Name,
		Password: hashedPassword,
		Status:   StatusMemberApplied,
	}, nil
}

func NewMemberEntityFromDoorayMember(doorayMember dtos.DoorayMember) MemberEntity {
	// 두레이 사용자의 경우 이미 두레이를 통해 인증된 사용자 이기 때문에 상태를 '승인' 설정
	return MemberEntity{
		Type:           TypeMemberDooray,
		DoorayId:       doorayMember.Id,
		DoorayUserCode: doorayMember.UserCode,
		Name:           doorayMember.Name,
		Status:         StatusMemberApproved,
	}
}

func NewMemberEntityFromGoogleMember(googleMember dtos.GoogleMember) MemberEntity {
	// 구글 워크스페이스 사용자의 경우 이미 구글 워크스페이스를 통해 인증된 사용자 이기 때문에 상태를 '승인' 설정
	return MemberEntity{
		Type:       TypeMemberGoogle,
		GoogleId:   googleMember.Id,
		GoogleMail: googleMember.Email,
		Name:       googleMember.Name,
		Picture:    googleMember.Picture,
		Status:     StatusMemberApproved,
	}
}
