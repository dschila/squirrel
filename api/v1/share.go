package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/proph/squirrel/services"
)

func ShareRoutes(r *gin.RouterGroup) {
	c, err := services.InitMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	shareRoute := r.Group("/share")
	shareRoute.POST("/", createFile(c.MinioClient))
}

func createFile(c *minio.Client) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		form, _ := ctx.MultipartForm()
		files := form.File["file"]

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
		}
		const msg = "files uploaded!"
		response := map[string]interface{}{
			"msg": msg,
		}

		ctx.JSON(http.StatusCreated, response)
	}
}
