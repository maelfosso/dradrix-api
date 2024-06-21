package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPreferencesCompany struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type UserPreferences struct {
	Company        UserPreferencesCompany `bson:"company" json:"company,omitempty"`
	OnboardingStep int                    `bson:"onboarding_step" json:"onboarding_step"`
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
