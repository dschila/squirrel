package database

import (
	"context"
	"time"

	"github.com/proph/squirrel/helpers"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB *mongo.Database
}

func (r *MongoDB) CLose() {
	logrus.Warning("Closing db connections")
}

func InitMongoDB(config helpers.Config) (*MongoDB, error) {
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(config.MONGO_URI))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &MongoDB{DB: mongoClient.Database(config.MONGO_DB)}, nil
}
