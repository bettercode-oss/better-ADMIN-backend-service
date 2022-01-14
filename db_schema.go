package main

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/member"
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
		&webhook.WebHookEntity{}, &webhook.WebHookMessageEntity{}); err != nil {
		return err
	}

	var permissionCount int64
	db.Raw("SELECT count(*) FROM permissions WHERE type= 'pre-define'").Scan(&permissionCount)

	if permissionCount == 0 {
		if err := db.Exec("INSERT INTO permissions(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", domain.PermissionManageSystemSettings, "시스템 설정(예. 두레이 로그인 등) 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", domain.PermissionManageMembers, "멤버 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", domain.PermissionManageAccessControl, "접근 제어 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", domain.PermissionManageOrganization, "조직 관리 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO permissions(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", domain.PermissionNoteWebHooks, "웹훅 전송 권한", time.Now(), time.Now()).Error; err != nil {
			return err
		}
	}

	var roleCount int64
	db.Raw("SELECT count(*) FROM roles WHERE type= 'pre-define'").Scan(&roleCount)

	if roleCount == 0 {
		if err := db.Exec("INSERT INTO roles(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", "시스템 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(1, 1)").Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO roles(type, name, description, created_at, updated_at) values(?, ?, ?, ?, ?)",
			"pre-define", "조직/멤버 관리자", "", time.Now(), time.Now()).Error; err != nil {
			return err
		}

		if err := db.Exec("INSERT INTO role_permissions(role_entity_id, permission_entity_id) values(2, 2),(2, 3),(2, 4)").Error; err != nil {
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

	return nil
}
