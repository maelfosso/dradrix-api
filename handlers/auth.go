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
	SavePinCode(ctx context.Context, pinCode models.PinCodeOTP) error
}

type GetPinCodeInput struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
}

func GetPinCodeFromPhoneNumber(mux chi.Router, a authInterface) {
	mux.Post("/auth/", func(w http.ResponseWriter, r *http.Request) {
		var input GetPinCodeInput

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
		pinCode := utils.GeneratePinCode(now)

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
		var m models.PinCodeOTP
		m.MessageId = res.Messages[0].ID
		m.PhoneNumber = input.PhoneNumber
		m.PinCode = pinCode

		a.SavePinCode(r.Context(), m)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func PinCodeVerification(mux chi.Router) {
	mux.Post("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		// read the request body
		// extract the phone number and the pin code
		// check that the pin code is 6 digit
		// check that the phone number is correct
		// check from the database that is the pin code associated with the phone number
		// if it's okay return okay
		// if not, return error message
	})
}
