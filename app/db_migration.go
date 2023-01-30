package app

import (
	"better-admin-backend-service/constants"
	memberDomain "better-admin-backend-service/member/domain"
	organizationDomain "better-admin-backend-service/organization/domain"
	rbacDomain "better-admin-backend-service/rbac/domain"
	siteDomain "better-admin-backend-service/site/domain"
	webhookDomain "better-admin-backend-service/webhook/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *App) migrateDatabase() error {
	log.Info(">>> Database Migrate")
	// 테이블 생성
	if err := a.gormDB.AutoMigrate(&memberDomain.MemberEntity{}, &siteDomain.SettingEntity{}, &rbacDomain.PermissionEntity{},
		&rbacDomain.RoleEntity{}, &organizationDomain.OrganizationEntity{},
		&webhookDomain.WebHookEntity{}, &webhookDomain.WebHookMessageEntity{}); err != nil {
		return err
	}

	var permissionCount int64
	a.gormDB.Raw("SELECT count(*) FROM permissions WHERE type= 'pre-define'").Scan(&permissionCount)

	if permissionCount == 0 {
		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			1, "pre-define", constants.PermissionManageSystemSettings, "시스템 설정(예. 두레이 로그인 등) 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", constants.PermissionManageMembers, "멤버 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			3, "pre-define", constants.PermissionManageAccessControl, "접근 제어 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			4, "pre-define", constants.PermissionManageOrganization, "조직 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			5, "pre-define", constants.PermissionNoteWebHooks, "웹훅 전송 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			6, "pre-define", constants.PermissionViewMonitoring, "모니터링 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}
	}

	var roleCount int64
	a.gormDB.Raw("SELECT count(*) FROM roles WHERE type= 'pre-define'").Scan(&roleCount)

	if roleCount == 0 {
		if err := a.gormDB.Exec("INSERT INTO roles(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			1, "pre-define", "시스템 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(1, 1), (1, 6)").Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO roles(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", "조직/멤버 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(2, 2),(2, 3),(2, 4)").Error; err != nil {
			return err
		}
	}

	// siteadm 계정 만들기
	var signId string
	a.gormDB.Raw("SELECT sign_id FROM members WHERE sign_id = ?", "siteadm").Scan(&signId)

	if len(signId) == 0 {
		if err := a.gormDB.Exec("INSERT INTO members(type, sign_id, name, password, status, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?)",
			"site", "siteadm", "사이트 관리자", "$2a$04$7Ca1ybGc4yFkcBnzK1C0qevHy/LSD7PuBbPQTZEs6tiNM4hAxSYiG", "approved", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		// 사이트 관리자에 사전 정의된 두가지 역할을 할당한다.(시스템 관리자, 멤버 관리자)
		if err := a.gormDB.Exec("INSERT INTO member_roles(member_entity_id, role_entity_id) values(1, 1),(1, 2)").Error; err != nil {
			return err
		}
	}

	return nil
}
