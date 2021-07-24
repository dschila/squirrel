package main

import (
	"log"
	"os"

	"github.com/proph/squirrel/api"
	"github.com/proph/squirrel/helpers"
)

func main() {
	// Load config
	config, err := helpers.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.Google.CredentialsJsonPath)

	api.InitServer(config)
}
