package services

import (
	"better-admin-backend-service/domain"
	"better-admin-backend-service/domain/webhook/entity"
	"better-admin-backend-service/domain/webhook/repository"
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
)

type WebHookService struct {
}

func (WebHookService) CreateWebHook(ctx context.Context, webHookInformation dtos.WebHookInformation) error {
	repository := repository.WebHookRepository{}
	lastEntity, err := repository.FindLast(ctx)
	var nextId uint
	if err != nil {
		if err == domain.ErrNotFound {
			nextId = 1
		} else {
			return err
		}
	} else {
		nextId = lastEntity.NextId()
	}

	entity, err := entity.NewWebHookEntity(ctx, nextId, webHookInformation)
	if err != nil {
		return err
	}

	return repository.Create(ctx, &entity)
}

func (WebHookService) GetWebHooks(ctx context.Context, pageable dtos.Pageable) ([]entity.WebHookEntity, int64, error) {
	return repository.WebHookRepository{}.FindAll(ctx, pageable)
}

func (WebHookService) DeleteWebHook(ctx context.Context, webHookId uint) error {
	repository := repository.WebHookRepository{}

	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	entity, err := repository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	entity.UpdatedBy = userClaim.Id
	return repository.Delete(ctx, entity)
}

func (WebHookService) GetWebHook(ctx context.Context, webHookId uint) (entity.WebHookEntity, error) {
	return repository.WebHookRepository{}.FindById(ctx, webHookId)
}

func (WebHookService) UpdateWebHook(ctx context.Context, webHookId uint, webHookInformation dtos.WebHookInformation) error {
	repository := repository.WebHookRepository{}

	entity, err := repository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	err = entity.Update(ctx, webHookInformation)
	if err != nil {
		return err
	}

	return repository.Save(ctx, entity)
}

func (WebHookService) NoteMessage(ctx context.Context, webHookId uint, message dtos.WebHookMessage) error {
	repository := repository.WebHookRepository{}
	entity, err := repository.FindById(ctx, webHookId)
	if err != nil {
		return err
	}

	entity.AddMessage(message)

	err = repository.Save(ctx, entity)
	if err != nil {
		return err
	}

	message.Title = entity.Name

	return entity.NoteMessage(message)
}
