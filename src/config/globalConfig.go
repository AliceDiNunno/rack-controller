package config

type GlobalConfig struct {
	CurrentEnvironment string
}

func LoadGlobalConfiguration() GlobalConfig {
	return GlobalConfig{
		CurrentEnvironment: RequireEnvString("ENV"),
	}
}
