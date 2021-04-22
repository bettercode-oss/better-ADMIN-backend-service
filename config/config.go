package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	JwtSecret string
	Dooray    struct {
		LdapDialUrl string
	}
}{}

func InitConfig(cfg string) {
	configor.Load(&Config, cfg)
}
