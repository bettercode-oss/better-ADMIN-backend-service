package config

import (
	"github.com/jinzhu/configor"
	"os"
)

const (
	EnvJwtSecret = "JWT_SECRET"
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

func InitConfig(file string) error {
	err := configor.Load(&Config, file)
	if err != nil {
		return err
	}

	if len(os.Getenv(EnvJwtSecret)) > 0 {
		Config.JwtSecret = os.Getenv(EnvJwtSecret)
	}

	return nil
}
