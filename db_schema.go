package main

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/logging"
	"better-admin-backend-service/domain/member"
	"better-admin-backend-service/domain/menu"
	"better-admin-backend-service/domain/organization"
	"better-admin-backend-service/domain/rbac"
	"better-admin-backend-service/domain/site"
	"better-admin-backend-service/domain/webhook"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func initializeDatabase(db *gorm.DB) error {
	fmt.Println(">>> InitializeDatabase")
	// 테이블 생성
	if err := db.AutoMigrate(&member.MemberEntity{}, &site.SettingEntity{}, &rbac.PermissionEntity{},
		&rbac.RoleEntity{}, &organization.OrganizationEntity{},
		&webhook.WebHookEntity{}, &webhook.WebHookMessageEntity{},
		menu.MenuEntity{}, &logging.MemberAccessLogEntity{}); err != nil {
		return err
	}

	var permissionCount int64
	db.Raw("SELECT count(*) FROM permissions WHERE type= 'pre-define'").Scan(&permissionCount)

	if permissionCount == 0 {
		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			1, "pre-define", domain.PermissionManageSystemSettings, "시스템 설정(예. 두레이 로그인 등) 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", domain.PermissionManageMembers, "멤버 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			3, "pre-define", domain.PermissionManageAccessControl, "접근 제어 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			4, "pre-define", domain.PermissionManageOrganization, "조직 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			5, "pre-define", domain.PermissionNoteWebHooks, "웹훅 전송 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			6, "pre-define", domain.PermissionManageMenus, "메뉴 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			7, "pre-define", domain.PermissionViewMonitoring, "모니터링 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}
	}

	var roleCount int64
	db.Raw("SELECT count(*) FROM roles WHERE type= 'pre-define'").Scan(&roleCount)

	if roleCount == 0 {
		if err := db.Exec("INSERT INTO roles(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			1, "pre-define", "시스템 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(1, 1), (1, 7)").Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO roles(id, type, name, description, created_at, updated_at, created_by, updated_by) values(?, ?, ?, ?, ?, ?, 1, 1)",
			2, "pre-define", "조직/멤버 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(2, 2),(2, 3),(2, 4),(2, 6)").Error; err != nil {
			return err
		}
	}

	var menuCount int64
	db.Raw("SELECT count(*) FROM menus WHERE deleted_at is NULL").Scan(&menuCount)

	if menuCount == 0 {
		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (1,?,?,NULL,'GNB1-22','ApartmentOutlined',NULL,0,NULL,1,1,0)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (2,?,?,NULL,'GNB2','BarcodeOutlined',NULL,0,NULL,1,1,1)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (3,?,?,NULL,'GNB3','ApartmentOutlined',NULL,1,NULL,1,1,2)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (4,?,?,NULL,'SNB1','ApiOutlined',NULL,0,1,1,1,0)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (5,?,?,NULL,'SNB2','CiCircleOutlined','/snb2',0,1,1,1,1)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (6,?,?,NULL,'Sub1','BranchesOutlined','/sub1',0,4,1,1,0)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (7,?,?,NULL,'Sub2','BulbOutlined','/sub2',0,4,1,1,1)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO menus (id,created_at,updated_at,deleted_at,name,icon,link,disabled,parent_menu_id,created_by,updated_by,sequence) "+
			"VALUES (8,?,?,NULL,'Sample','CarryOutOutlined','/sample-list',0,2,1,1,0)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}
	}

	// siteadm 계정 만들기
	var signId string
	db.Raw("SELECT sign_id FROM members WHERE sign_id = ?", "siteadm").Scan(&signId)

	if len(signId) == 0 {
		if err := db.Exec("INSERT INTO members(type, sign_id, name, password, status, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?)",
			"site", "siteadm", "사이트 관리자", "$2a$04$7Ca1ybGc4yFkcBnzK1C0qevHy/LSD7PuBbPQTZEs6tiNM4hAxSYiG", "approved", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		// 사이트 관리자에 사전 정의된 두가지 역할을 할당한다.(시스템 관리자, 멤버 관리자)
		if err := db.Exec("INSERT INTO member_roles(member_entity_id, role_entity_id) values(1, 1),(1, 2)").Error; err != nil {
			return err
		}
	}

	// 기본 settings
	var siteSettingCount int64
	db.Raw("SELECT count(*) FROM site_settings WHERE deleted_at is NULL").Scan(&siteSettingCount)

	if siteSettingCount == 0 {
		if err := db.Exec("INSERT INTO site_settings (id, created_at, updated_at, deleted_at, key, value, created_by, updated_by) "+
			"VALUES (1,?,?,NULL,'member-access-log','{\"retentionDays\":30}', 1, 1)",
			time.Now(), time.Now()).Error; err != nil {
			return err
		}
	}

	return nil
}
