package repository

import (
	"net/http"

	"github.com/proph/squirrel/database"
	"github.com/proph/squirrel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ShareEntity IShare

type shareEntity struct {
	database   *database.MongoDB
	repository *mongo.Collection
}

type IShare interface {
	CreateShare(shareEntity models.Share) (models.Share, int, error)
}

func NewShareEntity(database *database.MongoDB) IShare {
	repo := database.DB.Collection("share")
	ShareEntity = &shareEntity{database: database, repository: repo}
	return ShareEntity
}

func (share *shareEntity) CreateShare(newShare models.Share) (models.Share, int, error) {
	newShare.Id = primitive.NewObjectID()
	ctx, cancel := initContext()
	defer cancel()
	_, err := share.repository.InsertOne(ctx, newShare)
	if err != nil {
		return models.Share{}, http.StatusInternalServerError, err
	}
	return newShare, http.StatusOK, nil
}
