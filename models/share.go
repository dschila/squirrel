package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Share struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Files []Files            `bson:"files" json:"files"`
}

type Files struct {
	Filename   string `bson:"filename" json:"filename"`
	StorageUrl string `bson:"storageUrl" json:"storageUrl"`
}
