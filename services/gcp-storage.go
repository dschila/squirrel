package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
)

type GCPStorage struct {
	StorageClient *storage.Client
}

func InitGCPStorageClient() (*GCPStorage, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return &GCPStorage{StorageClient: client}, err
}

func CreateBucket(storage *GCPStorage, bucketName string) (bool, error) {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := viper.GetString("google.projectid")

	// Creates a Bucket instance.
	bucket := storage.StorageClient.Bucket(bucketName)

	// Creates the new bucket.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket %v created.\n", bucketName)
	return true, nil
}
