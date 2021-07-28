package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/models"
	"github.com/proph/squirrel/repository"
	"github.com/proph/squirrel/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ShareRoutes(r *gin.RouterGroup, mongo *database.MongoDB) {
	c, err := services.InitMinioClient()
	if err != nil {
		log.Fatal(err)
	}
	const bucket = "squirrel-bucket"
	_, err = services.MakeBucket(c, bucket)
	if err != nil {
		log.Fatal(err)
	}

	shareEntity := repository.NewShareEntity(mongo)

	shareRoute := r.Group("/share")
	shareRoute.POST("/", createFile(c.MinioClient, bucket, shareEntity))
	shareRoute.GET("/:shareId", getShare(shareEntity))
}

func getShare(shareEntity repository.IShare) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		shareId := ctx.Param("shareId")
		share, code, err := shareEntity.GetShare(shareId)
		if err != nil {
			logrus.Warn("Share not found.")
		}
		ctx.JSON(code, share)
	}
}

func createFile(c *minio.Client, bucket string, shareEntity repository.IShare) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_ = os.Mkdir("tmp", 0700)

		form, _ := ctx.MultipartForm()
		files := form.File["file"]

		var shareFiles []models.Files

		for _, file := range files {
			filePath := fmt.Sprintf("tmp/%s", file.Filename)
			log.Printf("Received: %s", filePath)
			err := ctx.SaveUploadedFile(file, filePath)
			if err != nil {
				log.Fatal(err)
			}

			contentType, err := mimetype.DetectFile(filePath)
			if err != nil {
				log.Fatal(err)
			}

			info, putError := c.FPutObject(ctx, bucket, file.Filename, filePath, minio.PutObjectOptions{ContentType: contentType.String()})
			if putError != nil {
				log.Fatal(putError)
			}
			log.Printf("Successfully uploaded %s of size %d\n", file.Filename, info.Size)

			err = os.Remove(filePath)
			if err != nil {
				log.Fatal(err)
			}

			shareFiles = append(shareFiles, models.Files{
				Id:       primitive.NewObjectID().Hex(),
				Filename: file.Filename,
				Path:     fmt.Sprintf("%s/%s", bucket, info.Key),
			})

		}

		newShare := models.Share{
			Files: shareFiles,
		}

		newShare, _, err := shareEntity.CreateShare(newShare)
		if err != nil {
			logrus.Error(err)
		}

		ctx.JSON(http.StatusCreated, newShare)
	}
}
