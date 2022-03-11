package logging

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"time"
)

const (
	TypeMemberPageAccessLog     = "PAGE_ACCESS"
	TypeMemberPageAccessLogName = "화면"
	TypeMemberApiAccessLog      = "API_ACCESS"
	TypeMemberApiAccessLogName  = "API"
)

type MemberAccessLogEntity struct {
	ID               uint    `gorm:"primarykey"`
	MemberId         uint    `gorm:"not null;index"`
	Type             string  `gorm:"type:varchar(20);not null"`
	Url              string  `gorm:"type:varchar(100);not null"`
	Method           *string `gorm:"type:varchar(10)"`
	Parameters       *string `gorm:"type:varchar(2000)"`
	Payload          *string `gorm:"type:varchar(2000)"`
	StatusCode       *uint
	IpAddress        string `gorm:"type:varchar(20)"`
	BrowserUserAgent string `gorm:"type:varchar(2000)"`
	CreatedAt        time.Time
}

func (MemberAccessLogEntity) TableName() string {
	return "member_access_logs"
}

func (m MemberAccessLogEntity) GetType() string {
	if TypeMemberApiAccessLog == m.Type {
		return TypeMemberApiAccessLogName
	}

	if TypeMemberPageAccessLog == m.Type {
		return TypeMemberPageAccessLogName
	}

	return ""
}

func NewMemberAccessLogEntity(ctx context.Context, accessLog dtos.MemberAccessLog) (MemberAccessLogEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return MemberAccessLogEntity{}, err
	}

	if accessLog.Type != TypeMemberPageAccessLog && accessLog.Type != TypeMemberApiAccessLog {
		return MemberAccessLogEntity{}, domain.ErrNotSupportedAccessLogType
	}

	return MemberAccessLogEntity{
		MemberId:         userClaim.Id,
		Type:             accessLog.Type,
		Url:              accessLog.Url,
		Method:           accessLog.Method,
		Parameters:       accessLog.Parameters,
		Payload:          accessLog.Payload,
		StatusCode:       accessLog.StatusCode,
		IpAddress:        accessLog.IpAddress,
		BrowserUserAgent: accessLog.GetHumanizeBrowserUserAgent(),
	}, nil
}
