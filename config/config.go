package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	JwtSecret string
	Dooray    struct {
		LdapDialUrl string
	}
	GoogleOAuth struct {
		OAuthUri string
		AuthUri  string
		TokenUri string
	}
}{}

func InitConfig(cfg string) {
	configor.Load(&Config, cfg)
}
