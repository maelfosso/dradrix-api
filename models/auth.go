package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber string `json:"phone_number,omitempty"`
	Name        string `json:"name,omitempty"`
}

type OTP struct {
	gorm.Model
	WaMessageId string `json:"wa_message_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PinCode     string `json:"pin_code,omitempty"`
	Active      bool   `json:"active,omitempty"`
}
