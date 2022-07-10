package config

import GoLoggerClient "github.com/AliceDiNunno/go-logger-client"

func LoadEventConfiguration() GoLoggerClient.ClientConfiguration {
	version, err := GetEnvString("LOGGER_APP_VERSION")
	if err != nil {
		version = "unknown"
	}

	env, err := GetEnvString("ENV")
	if err != nil {
		env = "unknown"
	}

	return GoLoggerClient.ClientConfiguration{
		ProjectId:   RequireEnvString("LOGGER_PROJECT_ID"),
		Key:         RequireEnvString("LOGGER_EVENT_KEY"),
		Environment: env,
		Version:     version,

		RemoveFieldsFromDebugOutput: false,
	}
}
