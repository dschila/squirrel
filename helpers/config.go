package helpers

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Minio    MinioConfiguration
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	Uri  string
	Name string
}

type MinioConfiguration struct {
	Host            string
	AccessKeyID     string
	SecretAccessKey string
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
	viper.SetDefault("minio.host", "localhost:6971")
	viper.SetDefault("minio.accesskeyid", "minioadmin")
	viper.SetDefault("minio.secretaccesskey", "minioadmin")

	err = viper.Unmarshal(&config)
	return
}
