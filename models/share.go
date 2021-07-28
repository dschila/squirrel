package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Share struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Files []Files            `bson:"files" json:"files"`
}

type Files struct {
	Id       string `bson:"id" json:"id"`
	Filename string `bson:"filename" json:"filename"`
	Path     string `bson:"path" json:"path"`
}
