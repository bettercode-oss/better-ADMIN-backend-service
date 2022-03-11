package dtos

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
	log "github.com/sirupsen/logrus"
	"time"
)

type MemberAccessLog struct {
	Id               uint      `json:"id"`
	MemberId         uint      `json:"memberId"`
	Type             string    `json:"type" validate:"required"`
	TypeName         string    `json:"typeName"`
	Url              string    `json:"url" validate:"required"`
	Method           *string   `json:"method,omitempty"`
	Parameters       *string   `json:"parameters,omitempty"`
	Payload          *string   `json:"payload,omitempty"`
	StatusCode       *uint     `json:"statusCode,omitempty"`
	IpAddress        string    `json:"ipAddress"`
	BrowserUserAgent string    `json:"browserUserAgent"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (m MemberAccessLog) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}

func (m MemberAccessLog) GetHumanizeBrowserUserAgent() string {
	ua := user_agent.New(m.BrowserUserAgent)

	humanizeUserAgent := map[string]interface{}{}
	humanizeUserAgent["mobile"] = ua.Mobile()
	humanizeUserAgent["platform"] = ua.Platform()
	humanizeUserAgent["os"] = ua.OS()

	name, version := ua.Engine()
	engine := map[string]interface{}{}
	engine["name"] = name
	engine["version"] = version
	humanizeUserAgent["engine"] = engine

	name, version = ua.Browser()
	browser := map[string]interface{}{}
	browser["name"] = name
	browser["version"] = version
	humanizeUserAgent["browser"] = browser

	jsonString, err := json.Marshal(humanizeUserAgent)
	if err != nil {
		log.Error(errors.New(err))
		return ""
	}

	return string(jsonString)
}
