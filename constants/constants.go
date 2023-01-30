package constants

const (
	// RBAC
	PreDefineTypeKey   = "pre-define"
	PreDefineTypeName  = "사전정의"
	UserDefineTypeKey  = "user-define"
	UserDefineTypeName = "사용자정의"

	// Permissaion
	PermissionManageAccessControl  = "MANAGE_ACCESS_CONTROL"
	PermissionManageMembers        = "MANAGE_MEMBERS"
	PermissionManageOrganization   = "MANAGE_ORGANIZATION"
	PermissionManageSystemSettings = "MANAGE_SYSTEM_SETTINGS"
	PermissionNoteWebHooks         = "NOTE_WEB_HOOKS"
	PermissionViewMonitoring       = "VIEW_MONITORING"

	// Member
	TypeMemberSite       = "site"
	TypeMemberSiteName   = "사이트"
	TypeMemberDooray     = "dooray"
	TypeMemberDoorayName = "두레이"
	TypeMemberGoogle     = "google"
	TypeMemberGoogleName = "구글"
	StatusMemberApplied  = "applied"
	StatusMemberApproved = "approved"

	// Settings
	SettingKeyDoorayLogin          = "dooray-login"
	SettingKeyGoogleWorkspaceLogin = "google-workspace-login"
	SettingKeyMemberAccessLog      = "member-access-log"
	SettingKeyAppVersion           = "app-version"
)
