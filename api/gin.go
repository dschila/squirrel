package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/proph/squirrel/api/v1"
	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/helpers"
	"github.com/sirupsen/logrus"
)

func InitServer(config helpers.Config) {
	database, err := database.InitMongoDB(config)
	if err != nil {
		logrus.Error(err)
	}

	server := gin.Default()
	server.Use(gin.Logger())

	route := server.Group("/api/v1")
	v1.ShareRoutes(route, database)

	server.Run(":" + config.PORT)
}
