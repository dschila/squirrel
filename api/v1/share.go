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
)

func ShareRoutes(r *gin.RouterGroup, mongo *database.MongoDB) {
	c, err := services.InitMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	_, err = services.MakeBucket(c, "orca")
	if err != nil {
		log.Fatal(err)
	}

	shareEntity := repository.NewShareEntity(mongo)

	shareRoute := r.Group("/share")
	shareRoute.POST("/", createFile(c.MinioClient, shareEntity))
}

func createFile(c *minio.Client, shareEntity repository.IShare) func(ctx *gin.Context) {
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

			info, putError := c.FPutObject(ctx, "orca", file.Filename, filePath, minio.PutObjectOptions{ContentType: contentType.String()})
			if putError != nil {
				log.Fatal(putError)
			}
			log.Printf("Successfully uploaded %s of size %d\n", file.Filename, info.Size)

			err = os.Remove(filePath)
			if err != nil {
				log.Fatal(err)
			}

			shareFiles = append(shareFiles, models.Files{
				Filename:   file.Filename,
				StorageUrl: info.Key,
			})
		}

		newShare := models.Share{
			Files: shareFiles,
		}

		newShare, _, err := shareEntity.CreateShare(newShare)
		if err != nil {
			logrus.Error(err)
		}

		response := map[string]interface{}{
			"share": newShare,
		}

		ctx.JSON(http.StatusCreated, response)
	}
}
