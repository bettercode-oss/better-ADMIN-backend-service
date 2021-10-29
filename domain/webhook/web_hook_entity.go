package webhook

import (
	"better-admin-backend-service/adapters"
	"better-admin-backend-service/domain"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/security"
	"gorm.io/gorm"
)

type WebHookEntity struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	AccessToken string
	Messages    []WebHookMessageEntity `gorm:"foreignKey:WebHookId"`
}

func (WebHookEntity) TableName() string {
	return "web_hooks"
}

type WebHookMessageEntity struct {
	gorm.Model
	WebHookId uint
	Message   string `gorm:"not null"`
}

func (WebHookMessageEntity) TableName() string {
	return "web_hook_messages"
}

func (w *WebHookEntity) Update(information dtos.WebHookInformation) {
	w.Name = information.Name
	w.Description = information.Description
}

func (w WebHookEntity) nextId() uint {
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

func NewWebHookEntity(id uint, information dtos.WebHookInformation) (WebHookEntity, error) {
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
	}, nil
}
