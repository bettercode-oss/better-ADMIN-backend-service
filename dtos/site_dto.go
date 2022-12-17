package dtos

import (
	"better-admin-backend-service/config"
	"fmt"
)

type DoorayLoginSetting struct {
	Used               *bool  `json:"used" binding:"required"`
	Domain             string `json:"domain" binding:"required_if=Used true"`
	AuthorizationToken string `json:"authorizationToken" binding:"required_if=Used true"`
}

type SiteSettingsSummary struct {
	DoorayLoginUsed          bool   `json:"doorayLoginUsed"`
	GoogleWorkspaceLoginUsed bool   `json:"googleWorkspaceLoginUsed"`
	GoogleWorkspaceOAuthUri  string `json:"googleWorkspaceOAuthUri"`
}

type GoogleWorkspaceLoginSetting struct {
	Used         *bool  `json:"used" binding:"required"`
	Domain       string `json:"domain" binding:"required_if=Used true"`
	ClientId     string `json:"clientId" binding:"required_if=Used true"`
	ClientSecret string `json:"clientSecret" binding:"required_if=Used true"`
	RedirectUri  string `json:"redirectUri" binding:"required_if=Used true"`
}

func (g GoogleWorkspaceLoginSetting) GetOAuthUri() string {
	return fmt.Sprintf("%v?client_id=%v&redirect_uri=%v&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&approval_prompt=force&access_type=offline",
		config.Config.GoogleOAuth.OAuthUri, g.ClientId, g.RedirectUri)
}

type AppVersionSetting struct {
	Version uint `json:"version"`
}

func (avs *AppVersionSetting) Increase() {
	avs.Version = avs.Version + 1
}

func NewAppVersionSetting() AppVersionSetting {
	return AppVersionSetting{
		Version: 1,
	}
}
