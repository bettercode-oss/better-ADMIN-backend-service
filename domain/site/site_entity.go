package site

import (
	"encoding/json"
	"gorm.io/gorm"
)

const (
	SettingKeyDoorayLogin = "dooray-login"
)

type SettingEntity struct {
	gorm.Model
	Key         string
	Value       string
	ValueObject interface{} `gorm:"-"`
}

func (SettingEntity) TableName() string {
	return "site_settings"
}

func (u *SettingEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b, err := json.Marshal(u.ValueObject)
	if err != nil {
		return
	}
	u.Value = string(b)
	return
}

func (u *SettingEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	b, err := json.Marshal(u.ValueObject)
	if err != nil {
		return
	}
	u.Value = string(b)
	return
}

func (u *SettingEntity) AfterFind(tx *gorm.DB) (err error) {
	err = json.Unmarshal([]byte(u.Value), &u.ValueObject)
	if err != nil {
		return
	}

	return
}
