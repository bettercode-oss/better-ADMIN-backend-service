package app

import (
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
			1, "pre-define", "access-control-permission.all", "권한 관리에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", "access-control-permission.create", "권한 생성", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			3, "pre-define", "access-control-permission.read", "권한 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			4, "pre-define", "access-control-permission.update", "권한 수정", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			5, "pre-define", "access-control-permission.delete", "권한 삭제", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			6, "pre-define", "access-control-role.all", "역할 관리에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			7, "pre-define", "access-control-role.create", "역할 생성", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			8, "pre-define", "access-control-role.read", "역할 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			9, "pre-define", "access-control-role.update", "역할 수정", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			10, "pre-define", "access-control-role.delete", "역할 삭제", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			11, "pre-define", "member.all", "멤버 관리에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			12, "pre-define", "member.read", "멤버 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			13, "pre-define", "member.update", "멤버 승인/거부 및 역할 할당", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			14, "pre-define", "organization.all", "조직 관리에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			15, "pre-define", "organization.create", "조직 생성", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			16, "pre-define", "organization.read", "조직 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			17, "pre-define", "organization.update", "조직 수정", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			18, "pre-define", "organization.delete", "조직 삭제", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			19, "pre-define", "site-settings.all", "사이트 셋팅에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			20, "pre-define", "site-settings.read", "사이트 셋팅 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			21, "pre-define", "site-settings.update", "사이트 셋팅 수정", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			22, "pre-define", "web-hook.all", "웹훅에 관한 모든 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			23, "pre-define", "web-hook.create", "웹훅 생성", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			24, "pre-define", "web-hook.read", "웹훅 조회", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			25, "pre-define", "web-hook.update", "웹훅 수정", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			26, "pre-define", "web-hook.delete", "웹훅 삭제", time.Now(), time.Now()).Error; err != nil {
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

		if err := a.gormDB.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(1, 1), (1, 6), (1, 19), (1, 22)").Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO roles(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", "조직/멤버 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := a.gormDB.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(2, 11),(2, 14)").Error; err != nil {
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
