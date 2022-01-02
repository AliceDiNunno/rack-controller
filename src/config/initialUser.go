package config

type InitialUserConfig struct {
	Mail        string
	Password    string
	AccessToken string
}

func LoadInitialUserConfiguration() *InitialUserConfig {
	requireUserConfiguration, _ := GetEnvBool("CONFIGURE_INITIAL_USER")

	if !requireUserConfiguration {
		return nil
	}

	mail := RequireEnvString("INITIAL_USER_MAIL")
	password := RequireEnvString("INITIAL_USER_PASSWORD")
	accessToken := RequireEnvString("INITIAL_USER_ACCESS_TOKEN")

	return &InitialUserConfig{
		Mail:        mail,
		Password:    password,
		AccessToken: accessToken,
	}
}
