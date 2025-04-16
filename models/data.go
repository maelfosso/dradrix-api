package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataAuthor struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type Data struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
	// key: code of the field
	// value: the value entered by the user
	// value type: depends on the type associated to the field when creating the activity
	Values map[string]any `bson:"values" json:"values"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at"`

	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id"`
	CreatedBy  DataAuthor         `bson:"created_by" json:"created_by"`
}

type UploadedFile struct {
	UploadedBy primitive.ObjectID `bson:"uploaded_by" json:"uploaded_by"`
	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id"`
	FileKey    string             `bson:"file_key" json:"file_key"`

	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}
