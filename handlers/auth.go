package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/requests"
	"stockinos.com/api/utils"
)

type authInterface interface {
	CreateOTP(ctx context.Context, pinCode models.OTP) error
	SaveOTP(ctx context.Context, pinCode models.OTP) error
	CheckOTP(ctx context.Context, phoneNumber, pinCode string) (models.OTP, error)
}

type GetOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
	PinCode     string `json:"pin_code,omitempty"`     // Pin code entered
}

func GetOTPFromPhoneNumber(mux chi.Router, a authInterface) {
	mux.Post("/auth/otp", func(w http.ResponseWriter, r *http.Request) {
		var input GetOTPRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate the pin code of 6 digits
		now := time.Now()
		pinCode := utils.GenerateOTP(now)

		// send the pin code to a the phone number using Whatsapp API
		res, err := requests.SendWoZOTP(
			input.PhoneNumber,
			input.Language,
			pinCode,
		)

		// if request failed, no whatsapp account with phone number
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// if not, save the association phone number/pin code in the db
		var m models.OTP
		m.WaMessageId = res.Messages[0].ID
		m.PhoneNumber = input.PhoneNumber
		m.PinCode = pinCode

		a.CreateOTP(r.Context(), m)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func OTPVerification(mux chi.Router, a authInterface) {
	mux.Post("/auth/verify-otp", func(w http.ResponseWriter, r *http.Request) {
		// read the request body
		var input VerifyOTPRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number and the pin code
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check that the pin code is 6 digit
		var m models.OTP

		// check that the phone number is correct
		m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		m.Active = false
		a.SaveOTP(r.Context(), m)
	})
}
