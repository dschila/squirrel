package services

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	MinioClient *minio.Client
}

func InitMinioClient() (*MinioStorage, error) {
	endpoint := "localhost:6970"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"

	credentails := credentials.NewStaticV4(accessKeyID, secretAccessKey, "")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentails,
		Secure: false,
	})

	return &MinioStorage{MinioClient: client}, err
}

func MakeBucket(storage *MinioStorage, bucketName string) (bool, error) {
	location := "de"
	ctx := context.Background()

	err := storage.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := storage.MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already on %s\n", bucketName)
			return true, nil
		} else {
			return false, err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return true, nil
}
