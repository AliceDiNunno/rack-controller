package config

type GinConfig struct {
	ListenAddress string
	Port          int
	TlsEnabled    bool
	Prefix        string
}

func LoadGinConfiguration() GinConfig {
	prefix, err := GetEnvString("HTTP_API_PREFIX")

	if err != nil {
		prefix = "/"
	}

	return GinConfig{
		ListenAddress: RequireEnvString("HTTP_LISTEN_ADDRESS"),
		Port:          RequireEnvInt("HTTP_LISTEN_PORT"),
		TlsEnabled:    RequireEnvBool("HTTP_TLS_ENABLED"),
		Prefix:        prefix,
	}
}
