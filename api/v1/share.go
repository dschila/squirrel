package v1

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/models"
	"github.com/proph/squirrel/repository"
	"github.com/sirupsen/logrus"
)

func ShareRoutes(r *gin.RouterGroup, mongo *database.MongoDB) {
	/*
		c, err := services.InitGCPStorageClient()
		if err != nil {
			log.Fatal(err)
		}
	*/
	/*	_, err = services.CreateBucket(c, "test-01")
		if err != nil {
			log.Fatal(err)
		}
	*/

	shareEntity := repository.NewShareEntity(mongo)

	shareRoute := r.Group("/share")
	shareRoute.POST("/", createFile(shareEntity))
}

func createFile(shareEntity repository.IShare) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()
		_ = os.Mkdir("tmp", 0700)

		form, _ := c.MultipartForm()
		files := form.File["file"]

		client, err := storage.NewClient(ctx)
		if err != nil {
			logrus.Error("storage.NewClient: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
		}
		defer client.Close()

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		var shareFiles []models.Files

		for _, file := range files {
			filePath := fmt.Sprintf("tmp/%s", file.Filename)
			logrus.Printf("Received: %s", filePath)
			err := c.SaveUploadedFile(file, filePath)
			if err != nil {
				log.Fatal(err)
			}

			//contentType, err := mimetype.DetectFile(filePath)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Open(filePath)
			if err != nil {
				logrus.Error("os.Open: %v", err)
			}
			defer f.Close()
			wc := client.Bucket("dsc-test-01").Object(file.Filename).NewWriter(ctx)
			if _, err = io.Copy(wc, f); err != nil {
				logrus.Error("io.Copy: %v", err)
			}
			if err := wc.Close(); err != nil {
				logrus.Error("Writer.Close: %v", err)
			}

			err = os.Remove(filePath)
			if err != nil {
				logrus.Fatal(err)
			}

			shareFiles = append(shareFiles, models.Files{
				Filename:   file.Filename,
				StorageUrl: "xxxxx",
			})
		}

		newShare := models.Share{
			Files: shareFiles,
		}

		newShare, _, err = shareEntity.CreateShare(newShare)
		if err != nil {
			logrus.Error(err)
		}

		response := map[string]interface{}{
			"share": newShare,
		}

		c.JSON(http.StatusCreated, response)
	}
}
