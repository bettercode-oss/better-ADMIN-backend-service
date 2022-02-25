package dtos

import (
	"github.com/labstack/echo"
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
	IpAddress        string    `json:"ipAddress"`
	BrowserUserAgent string    `json:"browserUserAgent"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (m MemberAccessLog) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}
