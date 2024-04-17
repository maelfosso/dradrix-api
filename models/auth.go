package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
}

type OTP struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	WaMessageId string             `bson:"wa_message_id" json:"wa_message_id,omitempty"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	PinCode     string             `bson:"pin_code" json:"pin_code,omitempty"`
	Active      bool               `bson:"active,omitempty" json:"active,omitempty"`
}
