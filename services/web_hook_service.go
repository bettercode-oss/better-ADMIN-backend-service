package services

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/errors"
	"better-admin-backend-service/helpers"
	"better-admin-backend-service/webhook/domain"
	"better-admin-backend-service/webhook/repository"
	"context"
)

type WebHookService struct {
	webHookRepository *repository.WebHookRepository
}

func NewWebHookService(webHookRepository *repository.WebHookRepository) *WebHookService {
	return &WebHookService{
		webHookRepository: webHookRepository,
	}
}

func (s WebHookService) CreateWebHook(ctx context.Context, webHookInformation dtos.WebHookInformation) error {
	lastEntity, err := s.webHookRepository.FindLast(ctx)
	var nextId uint
	if err != nil {
		if err == errors.ErrNotFound {
			nextId = 1
		} else {
			return err
		}
	} else {
		nextId = lastEntity.NextId()
	}

	entity, err := domain.NewWebHookEntity(ctx, nextId, webHookInformation)
	if err != nil {
		return err
	}

	return s.webHookRepository.Create(ctx, &entity)
}

func (s WebHookService) GetWebHooks(ctx context.Context, pageable dtos.Pageable) ([]domain.WebHookEntity, int64, error) {
	return s.webHookRepository.FindAll(ctx, pageable)
}

func (s WebHookService) DeleteWebHook(ctx context.Context, webHookId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	entity, err := s.webHookRepository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	entity.UpdatedBy = userClaim.Id
	return s.webHookRepository.Delete(ctx, entity)
}

func (s WebHookService) GetWebHook(ctx context.Context, webHookId uint) (domain.WebHookEntity, error) {
	return s.webHookRepository.FindById(ctx, webHookId)
}

func (s WebHookService) UpdateWebHook(ctx context.Context, webHookId uint, webHookInformation dtos.WebHookInformation) error {
	entity, err := s.webHookRepository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	err = entity.Update(ctx, webHookInformation)
	if err != nil {
		return err
	}

	return s.webHookRepository.Save(ctx, entity)
}

func (s WebHookService) NoteMessage(ctx context.Context, webHookId uint, message dtos.WebHookMessage) error {
	entity, err := s.webHookRepository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	entity.AddMessage(message)

	err = s.webHookRepository.Save(ctx, entity)
	if err != nil {
		return err
	}

	message.Title = entity.Name

	return entity.NoteMessage(message)
}
