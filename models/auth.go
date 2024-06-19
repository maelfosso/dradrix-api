package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Preferences struct {
	Company struct {
		Id   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
		Name string             `bson:"name" json:"name,omitempty"`
	} `bson:"company" json:"company,omitempty"`
	CurrentOnboardingStep int `bson:"current_onboarding_step" json:"current_onboarding_step,omitempty"`
}
type User struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`

	Preferences Preferences `bson:"preferences" json:"preferences,omitempty"`
}

type OTP struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	WaMessageId string             `bson:"wa_message_id" json:"wa_message_id,omitempty"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	PinCode     string             `bson:"pin_code" json:"pin_code,omitempty"`
	Active      bool               `bson:"active,omitempty" json:"active,omitempty"`
}
