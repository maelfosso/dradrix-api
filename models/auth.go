package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPreferencesCompany struct {
	Id   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string             `bson:"name" json:"name,omitempty"`
}

type UserPreferences struct {
	Company               UserPreferencesCompany `bson:"company" json:"company,omitempty"`
	CurrentOnboardingStep int                    `bson:"current_onboarding_step" json:"current_onboarding_step,omitempty"`
}

type User struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	FirstName   string             `bson:"first_name,omitempty" json:"last_name,omitempty"`
	LastName    string             `bson:"last_name,omitempty" json:"first_name,omitempty"`

	Preferences UserPreferences `bson:"preferences" json:"preferences,omitempty"`
}

type OTP struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	WaMessageId string             `bson:"wa_message_id" json:"wa_message_id,omitempty"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number,omitempty"`
	PinCode     string             `bson:"pin_code" json:"pin_code,omitempty"`
	Active      bool               `bson:"active,omitempty" json:"active,omitempty"`
}
