package config

type GlobalConfig struct {
	CurrentEnvironment        string
	DebugEnvironmentVariables bool
}

func LoadGlobalConfiguration() GlobalConfig {
	return GlobalConfig{
		CurrentEnvironment:        RequireEnvString("ENV"),
		DebugEnvironmentVariables: RequireEnvBool("DEBUG_ENV_VARIABLES"),
	}
}
