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
		Host:     RequireEnvString("DB_HOST"),
		Port:     RequireEnvInt("DB_PORT"),
		User:     RequireEnvString("POSTGRES_USER"),
		Password: RequireEnvString("POSTGRES_USER"),
		DbName:   RequireEnvString("POSTGRES_DB"),
	}
}
