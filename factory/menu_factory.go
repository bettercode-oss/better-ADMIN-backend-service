package factory

import (
	"better-admin-backend-service/domain/menu"
	"better-admin-backend-service/dtos"
)

func NewMenuInformationFromEntity(entity menu.MenuEntity) dtos.MenuInformation {
	menusAccessPermission := make([]dtos.MenuAccessPermission, 0)
	for _, permission := range entity.Permissions {
		menusAccessPermission = append(menusAccessPermission, dtos.MenuAccessPermission{
			Id:   permission.ID,
			Name: permission.Name,
		})
	}

	menuInformation := dtos.MenuInformation{
		Id:                entity.ID,
		Name:              entity.Name,
		Title:             entity.Name,
		Icon:              entity.Icon,
		Link:              entity.Link,
		Disabled:          entity.Disabled,
		AccessPermissions: menusAccessPermission,
	}

	menuInformation.SetUpMenuType()

	return menuInformation
}
