package services

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type MinioStorage struct {
	MinioClient *minio.Client
}

func InitMinioClient() (*MinioStorage, error) {
	endpoint := viper.GetString("minio.host")
	accessKeyID := viper.GetString("minio.accesskeyid")
	secretAccessKey := viper.GetString("minio.secretaccesskey")

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

	// when new bucket created set access policy to public
	policy := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`, bucketName)
	err = storage.MinioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return false, err
	} else {
		log.Printf("Successfully set access policy to public \n")
	}

	return true, nil
}
