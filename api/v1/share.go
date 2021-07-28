package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/proph/squirrel/controllers"
	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/repository"
)

func ShareRoutes(r *gin.RouterGroup, mongo *database.MongoDB) {
	shareEntity := repository.NewShareEntity(mongo)

	shareRoute := r.Group("/share")
	shareRoute.POST("/", controllers.UploadFileToBucket(shareEntity))
}
