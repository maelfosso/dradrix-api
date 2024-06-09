package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at,omitempty"`

	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}
