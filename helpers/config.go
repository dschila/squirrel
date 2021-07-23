package helpers

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfigurations

	// GCP
	GOOGLE_APPLICATION_CREDENTIALS string
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfigurations struct {
	DBUri  string
	DBName string
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
