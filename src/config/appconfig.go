package config

var appConfig *AppConfig

type AppConfig struct {
	UTMPath string

	HTTPPort uint16
	HTTPHost string
}

func NewAppConfig() *AppConfig {
	appConfig := &AppConfig{}
	appConfig.HTTPPort = 8788

	return appConfig
}

func GetAppConfig() *AppConfig {
	if appConfig == nil {
		panic("AppConfig is not initialized")
	}

	return appConfig
}
