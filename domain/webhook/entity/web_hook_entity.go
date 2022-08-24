package entity

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/security"
	"context"
	"gorm.io/gorm"
)

type WebHookEntity struct {
	gorm.Model
	Name        string                 `gorm:"type:varchar(100);not null"`
	Description string                 `gorm:"type:varchar(1000)"`
	AccessToken string                 `gorm:"type:varchar(1000)"`
	Messages    []WebHookMessageEntity `gorm:"foreignKey:WebHookId"`
	CreatedBy   uint
	UpdatedBy   uint
}

func (WebHookEntity) TableName() string {
	return "web_hooks"
}

type WebHookMessageEntity struct {
	gorm.Model
	WebHookId uint   `gorm:"not null"`
	Message   string `gorm:"type:text;not null"`
}

func (WebHookMessageEntity) TableName() string {
	return "web_hook_messages"
}

func (w *WebHookEntity) Update(ctx context.Context, information dtos.WebHookInformation) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	w.Name = information.Name
	w.Description = information.Description
	w.UpdatedBy = userClaim.Id

	return nil
}

func (w WebHookEntity) NextId() uint {
	return w.ID + 1
}

func (w *WebHookEntity) AddMessage(message dtos.WebHookMessage) {
	w.Messages = append(w.Messages, WebHookMessageEntity{Message: message.Text})
}

func (w WebHookEntity) NoteMessage(message dtos.WebHookMessage) error {
	if err := adapters.WebSocketAdapter().BroadcastMessage(message); err != nil {
		return err
	}

	return nil
}

func NewWebHookEntity(ctx context.Context, id uint, information dtos.WebHookInformation) (WebHookEntity, error) {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return WebHookEntity{}, err
	}

	accessToken, err := security.JwtAuthentication{}.GenerateJwtAccessTokenNeverExpired(security.UserClaim{
		Id:          id,
		Permissions: []string{domain.PermissionNoteWebHooks},
	})

	if err != nil {
		return WebHookEntity{}, nil
	}

	return WebHookEntity{
		Name:        information.Name,
		Description: information.Description,
		AccessToken: accessToken,
		CreatedBy:   userClaim.Id,
		UpdatedBy:   userClaim.Id,
	}, nil
}
