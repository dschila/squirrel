package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/proph/squirrel/models"
	"github.com/proph/squirrel/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func UploadFileToBucket(shareEntity repository.IShare) func(c *gin.Context) {
	return func(c *gin.Context) {
		_ = os.Mkdir("tmp", 0700)

		form, _ := c.MultipartForm()
		files := form.File["file"]

		client, err := storage.NewClient(c, option.WithCredentialsFile(viper.GetString("google.credentialsjsonpath")))
		if err != nil {
			logrus.Error("storage.NewClient: %v", err)
			c.JSON(http.StatusInternalServerError, nil)
		}
		defer client.Close()

		var shareFiles []models.Files
		logrus.Info("shareFiles", shareFiles)
		for _, file := range files {
			filePath := fmt.Sprintf("tmp/%s", file.Filename)
			logrus.Info("Received: %s", filePath)
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
			wc := client.Bucket(viper.GetString("google.bucketid")).Object(file.Filename).NewWriter(c)
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
				StorageUrl: "",
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
