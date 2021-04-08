package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	JwtSecret string
}{}

func InitConfig(cfg string) {
	configor.Load(&Config, cfg)
}
