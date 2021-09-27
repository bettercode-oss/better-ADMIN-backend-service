package dtos

import (
	"better-admin-backend-service/config"
	"errors"
	"fmt"
	"github.com/labstack/echo"
)

type DoorayLoginSetting struct {
	Used               *bool  `json:"used" validate:"required"`
	Domain             string `json:"domain"`
	AuthorizationToken string `json:"authorizationToken"`
}

func (d DoorayLoginSetting) Validate(ctx echo.Context) error {
	if err := ctx.Validate(d); err != nil {
		return err
	}

	if *d.Used == true {
		if len(d.Domain) == 0 || len(d.AuthorizationToken) == 0 {
			return errors.New("domain and authorizationToken are required")
		}
	}

	return nil
}

type SiteSettingsSummary struct {
	DoorayLoginUsed          bool   `json:"doorayLoginUsed"`
	GoogleWorkspaceLoginUsed bool   `json:"googleWorkspaceLoginUsed"`
	GoogleWorkspaceOAuthUri  string `json:"googleWorkspaceOAuthUri"`
}

type GoogleWorkspaceLoginSetting struct {
	Used         *bool  `json:"used" validate:"required"`
	Domain       string `json:"domain"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectUri  string `json:"redirectUri"`
}

func (g GoogleWorkspaceLoginSetting) Validate(ctx echo.Context) error {
	if err := ctx.Validate(g); err != nil {
		return err
	}

	if *g.Used == true {
		if len(g.Domain) == 0 || len(g.ClientId) == 0 || len(g.ClientSecret) == 0 || len(g.RedirectUri) == 0 {
			return errors.New("domain, clientId, clientSecret and redirectUri are required")
		}
	}

	return nil
}

func (g GoogleWorkspaceLoginSetting) GetOAuthUri() string {
	return fmt.Sprintf("%v?client_id=%v&redirect_uri=%v&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&approval_prompt=force&access_type=offline",
		config.Config.GoogleOAuth.OAuthUri, g.ClientId, g.RedirectUri)
}
