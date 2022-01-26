package config

type MongodbConfig struct {
	Host string
	Port int
}

func LoadMongodbConfiguration() MongodbConfig {
	return MongodbConfig{
		Host: RequireEnvString("MONGO_HOST"),
		Port: RequireEnvInt("MONGO_PORT"),
	}
}
