package config

type GormConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func LoadGormConfiguration() GormConfig {
	return GormConfig{
		Host:     RequireEnvString("POSTGRES_HOST"),
		Port:     RequireEnvInt("POSTGRES_PORT"),
		User:     RequireEnvString("POSTGRES_USER"),
		Password: RequireEnvString("POSTGRES_PASSWORD"),
		DbName:   RequireEnvString("POSTGRES_DB"),
	}
}
