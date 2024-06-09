package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityFieldsOptions struct {
	Multiple     bool    `bson:"multiple" json:"multiple,omitempty"` // List of values
	Automatic    bool    `bson:"automatic" json:"automatic,omitempty"`
	DefaultValue *string `bson:"default_value" json:"default_value,omitempty"`
	Reference    *string `bson:"reference" json:"reference,omitempty"` // Is it an id from another activity
}

type ActivityFields struct {
	Name        string                `bson:"name" json:"name,omitempty"`
	Description string                `bson:"description" json:"description,omitempty"`
	Type        string                `bson:"type" json:"type,omitempty"`       // Text, Number, Date, Time, Uploaded file
	Id          bool                  `bons:"id" json:"id,omitempty"`           // Is it an identifiant?
	Options     ActivityFieldsOptions `bson:"options" json:"options,omitempty"` // There can be options
	Code        string                `bson:"code" json:"code,omitempty"`       // the id associated to the field, created internally
}

type Activity struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Fields      []ActivityFields   `bson:"fields" json:"fields,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at,omitempty"`

	CompanyId primitive.ObjectID `bson:"company_id" json:"company_id,omitempty"`
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}

type Data struct {
	Id primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	// key: code of the field
	// value: the value entered by the user
	// value type: depends on the type associated to the field when creating the activity
	Values map[string]interface{} `bson:"values" json:"values,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at,omitempty"`

	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}
