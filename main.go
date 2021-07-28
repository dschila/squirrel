package main

import (
	"github.com/proph/squirrel/api"
	"github.com/proph/squirrel/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load config
	config, err := helpers.LoadConfig()
	if err != nil {
		logrus.Warn("No configuration file found.")
	}

	api.InitServer(config)
}
