package config

type OVHConfig struct {
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
}

func LoadOvhConfiguration() OVHConfig {
	return OVHConfig{
		Endpoint:          RequireEnvString("OVH_ENDPOINT"),
		ApplicationKey:    RequireEnvString("OVH_APPLICATION_KEY"),
		ApplicationSecret: RequireEnvString("OVH_APPLICATION_SECRET"),
		ConsumerKey:       RequireEnvString("OVH_CONSUMER_KEY"),
	}
}
