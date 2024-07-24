package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityFieldOptions struct {
	Multiple     bool    `bson:"multiple" json:"multiple"` // List of values
	Automatic    bool    `bson:"automatic" json:"automatic"`
	DefaultValue *string `bson:"default_value" json:"default_value"`
	Reference    *string `bson:"reference" json:"reference"` // Is it an id from another activity
}

type ActivityFieldUpload struct {
	TypeOfFiles      []string `bson:"type_of_files" json:"type_of_files,omitempty"`
	MaxNumberOfFiles int      `bson:"max_number_of_items" json:"max_number_of_items"`
}

type ActivityFieldMultipleChoices struct {
	Multiple bool     `bson:"multiple" json:"multiple,omitempty"`
	Choices  []string `bson:"choices" json:"choices,omitempty"`
}

type ActivityFieldKey struct {
	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id,omitempty"`
	FieldId    primitive.ObjectID `bson:"field_id" json:"field_id,omitempty"`
}

type ActivityFieldType struct {
	*ActivityFieldMultipleChoices `bson:",inline" json:",inline"`
	*ActivityFieldKey             `bson:",inline" json:",inline"`
	*ActivityFieldUpload          `bson:",inline" json:",inline"`
}

func NewActivityFieldType(fieldType string) ActivityFieldType {
	switch fieldType {
	case "multiple-choices":
		return ActivityFieldType{
			&ActivityFieldMultipleChoices{
				Multiple: false,
				Choices:  []string{},
			},
			nil,
			nil,
		}
	case "key":
		return ActivityFieldType{
			nil,
			&ActivityFieldKey{
				ActivityId: primitive.NilObjectID,
				FieldId:    primitive.NilObjectID,
			},
			nil,
		}
	case "upload":
		return ActivityFieldType{
			nil,
			nil,
			&ActivityFieldUpload{
				TypeOfFiles:      []string{},
				MaxNumberOfFiles: 0,
			},
		}

	default:
		return ActivityFieldType{
			ActivityFieldMultipleChoices: nil,
			ActivityFieldKey:             nil,
			ActivityFieldUpload:          nil,
		}
	}
}

type ActivityField struct {
	Id          primitive.ObjectID   `bson:"_id" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Type        string               `bson:"type" json:"type"`       // Text, Number, Date, Time, Uploaded file
	PrimaryKey  bool                 `bons:"key" json:"primary_key"` // Is it an identifier?
	Options     ActivityFieldOptions `bson:"options" json:"options"` // There can be options
	Code        string               `bson:"code" json:"code"`       // the id associated to the field, created internally
	Details     ActivityFieldType    `bson:"details" json:"details"`
}

type Activity struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Fields      []ActivityField    `bson:"fields" json:"fields"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at"`

	OrganizationId primitive.ObjectID `bson:"organization_id" json:"organization_id"`
	CreatedBy      primitive.ObjectID `bson:"created_by" json:"created_by"`
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
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by"`
}
