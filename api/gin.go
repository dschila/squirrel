package api

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/proph/squirrel/api/v1"
	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/helpers"
	"github.com/sirupsen/logrus"
)

func InitServer(config helpers.Configuration) {
	database, err := database.InitMongoDB(config)
	if err != nil {
		logrus.Error(err)
	}

	server := gin.Default()
	server.Use(gin.Logger())
	server.Use(cors.Default())
	route := server.Group("/api/v1")
	v1.ShareRoutes(route, database)

	server.Run(":" + strconv.Itoa(config.Server.Port))
}
