package models

import (
	"math"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	// "golang.org/x/exp/constraints"
)

type ActivityFieldText struct {
}

func (f ActivityFieldText) IsValid(value any) (string, bool) {
	return value.(string), true
}

type ActivityFieldDate struct {
	Layout string // Example of date
}

func (f ActivityFieldDate) IsValid(value any) (time.Time, bool) {
	d, err := time.Parse(f.Layout, value.(string))
	return d, err == nil
}

type ActivityFieldTime struct {
	Layout string // 12h AM/PM or 24h
}

func (f ActivityFieldTime) IsValid(value any) (time.Time, bool) {
	t, err := time.Parse(f.Layout, value.(string))
	return t, err == nil
}

type ActivityFieldNumber struct {
	Minimum *float64
	Maximum *float64
}

func (f ActivityFieldNumber) IsValid(value any) (float64, bool) {
	v, err := strconv.ParseFloat(value.(string), 32)
	if err != nil {
		return math.NaN(), false
	}

	if f.Maximum != nil && v > *f.Maximum {
		return math.NaN(), false
	}
	if f.Minimum != nil && v < *f.Minimum {
		return math.NaN(), false
	}

	return v, true
}

type ActivityFieldType interface {
	*ActivityFieldDate | *ActivityFieldTime | *ActivityFieldText | *ActivityFieldNumber
	IsValid(value any) (interface{}, bool)
}

type ActivityField[T ActivityFieldType] struct {
	Code        string               `bson:"code" json:"code,omitempty"` // the id associated to the field, created internally
	Name        string               `bson:"name" json:"name,omitempty"`
	Description string               `bson:"description" json:"description,omitempty"`
	Type        string               `bson:"type" json:"type,omitempty"`       // Text, Number, Date, Time, Uploaded file
	Id          bool                 `bons:"id" json:"id,omitempty"`           // Is it an identifiant?
	Options     ActivityFieldOptions `bson:"options" json:"options,omitempty"` // There can be options

	Details T
}

func NewActivityField[T ActivityFieldType](field T) *ActivityField[T] {
	return &ActivityField[T]{
		Details: field,
	}
}

type ActivityFieldOptions struct {
	Multiple     bool    `bson:"multiple" json:"multiple,omitempty"` // List of values
	Automatic    bool    `bson:"automatic" json:"automatic,omitempty"`
	DefaultValue *string `bson:"default_value" json:"default_value,omitempty"`
	Reference    *string `bson:"reference" json:"reference,omitempty"` // Is it an id from another activity
}

type Activity struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Fields      []ActivityField    `bson:"fields" json:"fields,omitempty"`

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
	Values map[string]any `bson:"values" json:"values,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at,omitempty"`

	ActivityId primitive.ObjectID `bson:"activity_id" json:"activity_id,omitempty"`
	CreatedBy  primitive.ObjectID `bson:"created_by" json:"created_by,omitempty"`
}
