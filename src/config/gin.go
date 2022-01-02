package config

type GinConfig struct {
	Host   string
	Port   int
	Mode   string
	Tls    bool
	Prefix string
}

func LoadGinConfiguration() GinConfig {
	prefix, err := GetEnvString("GIN_PREFIX")

	if err != nil {
		prefix = "/"
	}

	return GinConfig{
		Host:   RequireEnvString("GIN_LISTEN_URL"),
		Port:   RequireEnvInt("GIN_LISTEN_PORT"),
		Mode:   RequireEnvString("GIN_MODE"),
		Tls:    RequireEnvBool("GIN_TLS"),
		Prefix: prefix,
	}
}
