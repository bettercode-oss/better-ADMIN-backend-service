package menu

import (
	"better-admin-backend-service/dtos"
	"better-admin-backend-service/helpers"
	"context"
	"github.com/wesovilabs/koazee"
	"strings"
)

type MenuService struct {
}

func (MenuService) CreateMenu(ctx context.Context, menuInformation dtos.MenuInformation) error {
	menuEntity, err := NewMenuEntity(ctx, menuInformation)
	if err != nil {
		return err
	}
	return menuRepository{}.Create(ctx, menuEntity)
}

func (MenuService) GetAllMenus(ctx context.Context) ([]MenuEntity, error) {
	entities, err := menuRepository{}.FindAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(entities); i++ {
		entities[i].generatePath(entities)
	}

	entitiesSortedByPath := koazee.StreamOf(entities).Sort(func(a, b MenuEntity) int {
		return strings.Compare(a.Path, b.Path)
	}).Out().Val().([]MenuEntity)

	return entitiesSortedByPath, nil
}

func (service MenuService) ChangePosition(ctx context.Context, menuId uint, position dtos.MenuPosition) error {
	repository := menuRepository{}
	entity, err := repository.FindById(ctx, menuId)
	if err != nil {
		return err
	}

	err = entity.ChangePosition(ctx, position)
	if err != nil {
		return err
	}

	err = repository.Save(ctx, &entity)
	if err != nil {
		return err
	}

	return service.reSequenceMenus(ctx, position)
}

func (MenuService) reSequenceMenus(ctx context.Context, position dtos.MenuPosition) error {
	repository := menuRepository{}

	filters := map[string]interface{}{}
	filters["parentMenuId"] = position.ParentMenuId

	menuEntities, err := repository.FindAll(ctx, filters)
	if err != nil {
		return err
	}

	for sequence, menuId := range position.SameDepthMenusSequence {
		for _, menuEntity := range menuEntities {
			if menuEntity.ID == menuId {
				menuEntity.Sequence = uint(sequence)
				err := repository.Save(ctx, &menuEntity)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

func (MenuService) DeleteMenu(ctx context.Context, menuId uint) error {
	userClaim, err := helpers.ContextHelper().GetUserClaim(ctx)
	if err != nil {
		return err
	}

	repository := menuRepository{}
	entity, err := repository.FindById(ctx, menuId)
	if err != nil {
		return err
	}

	childEntities, err := entity.FindChildEntities(ctx)
	if err != nil {
		return err
	}

	for _, childEntity := range childEntities {
		childEntity.UpdatedBy = userClaim.Id
		if err := repository.Delete(ctx, childEntity); err != nil {
			return err
		}
	}

	entity.UpdatedBy = userClaim.Id
	return repository.Delete(ctx, entity)
}

func (MenuService) UpdateMenu(ctx context.Context, menuId uint, menuInformation dtos.MenuInformation) error {
	repository := menuRepository{}
	entity, err := repository.FindById(ctx, menuId)
	if err != nil {
		return err
	}

	if err := entity.Update(ctx, menuInformation); err != nil {
		return err
	}

	return repository.Save(ctx, &entity)
}
