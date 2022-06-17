package config

import GoLoggerClient "github.com/AliceDiNunno/go-logger-client"

func LoadEventConfiguration() GoLoggerClient.ClientConfiguration {
	version, err := GetEnvString("APP_VERSION")
	if err != nil {
		version = "unknown"
	}

	env, err := GetEnvString("ENV")
	if err != nil {
		env = "unknown"
	}

	return GoLoggerClient.ClientConfiguration{
		ProjectId:   RequireEnvString("PROJECT_ID"),
		Key:         RequireEnvString("EVENT_KEY"),
		Environment: env,
		Version:     version,

		RemoveFieldsFromDebugOutput: false,
	}
}
