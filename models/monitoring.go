package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityFields struct {
	Name        string `bson:"name" json:"name,omitempty"`
	Description string `bson:"description" json:"description,omitempty"`
	Type        string `bson:"type" json:"type,omitempty"`
}

type Activity struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Fields      []ActivityFields   `bson:"fields" json:"fields,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt   *time.Time         `bson:"deleted_at" json:"deleted_at,omitempty"`

	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}

type Data struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time         `bson:"deleted_at" json:"deleted_at,omitempty"`

	Values map[string]interface{} `bson:"values" json:"values,omitempty"`

	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}
