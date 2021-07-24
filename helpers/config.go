package helpers

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfigurations
	Google   GoogleCloudConfiguration
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfigurations struct {
	DBUri  string
	DBName string
}

type GoogleCloudConfiguration struct {
	ProjectID           string
	CredentialsJsonPath string
}

func LoadConfig() (config Configuration, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Set undefined variables
	viper.SetDefault("server.port", 6970)

	err = viper.Unmarshal(&config)
	return
}
