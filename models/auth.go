package models

import "gorm.io/gorm"

type PinCodeOTP struct {
	gorm.Model
	MessageId   string `json:"message_id",omitempty`
	PhoneNumber string `json:"phone_number",omitempty`
	PinCode     string `json:"pin_code",omitempty`
}
