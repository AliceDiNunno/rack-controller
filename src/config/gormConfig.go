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
		User:     RequireEnvString("DB_USER"),
		Password: RequireEnvString("DB_PASSWORD"),
		DbName:   RequireEnvString("DB_NAME"),
	}
}
