package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/proph/squirrel/api/v1"
)

func InitServer() {
	server := gin.Default()
	server.Use(gin.Logger())

	route := server.Group("/api/v1")
	v1.ShareRoutes(route)

	server.Run(":6970")
}
