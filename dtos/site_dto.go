package dtos

import (
	"errors"
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
	DoorayLoginUsed bool `json:"doorayLoginUsed"`
}
