package config

import (
	"log/slog"

	"github.com/Netflix/go-env"
)

var appConfig *AppConfig

type AppConfig struct {
	UTMPath string `env:"UTM_PATH,default=/Applications/UTM.app"`

	HTTPPort uint16 `env:"HTTP_PORT,default=8080"`
	HTTPHost string `env:"HTTP_HOST,default=0.0.0.0"`
}

func Bootstrap(l *slog.Logger) *AppConfig {
	appConfig = &AppConfig{}

	_, err := env.UnmarshalFromEnviron(appConfig)
	if err != nil {
		l.Error("Failed to load app config", "err", err)
		panic("Failed to load app config")
	}

	return appConfig
}

func GetAppConfig() *AppConfig {
	if appConfig == nil {
		panic("AppConfig is not initialized")
	}

	return appConfig
}
