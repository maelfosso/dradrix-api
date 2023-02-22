package models

import "gorm.io/gorm"

type OTP struct {
	gorm.Model
	WaMessageId string `json:"wa_message_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	PinCode     string `json:"pin_code,omitempty"`
	Active      bool   `json:"active,omitempty"`
}
