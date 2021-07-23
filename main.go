package main

import (
	"log"

	"github.com/proph/squirrel/api"
	"github.com/proph/squirrel/helpers"
)

func main() {
	// Load config
	config, err := helpers.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	api.InitServer(config)
}
