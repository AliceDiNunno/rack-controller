package config

type GinConfig struct {
	ListenAddress string
	Port          int
	TlsEnabled    bool
	Prefix        string
}

func LoadGinConfiguration() GinConfig {
	prefix, err := GetEnvString("API_PREFIX")

	if err != nil {
		prefix = "/"
	}

	return GinConfig{
		ListenAddress: RequireEnvString("LISTEN_ADDRESS"),
		Port:          RequireEnvInt("PORT"),
		TlsEnabled:    RequireEnvBool("TLS_ENABLED"),
		Prefix:        prefix,
	}
}
