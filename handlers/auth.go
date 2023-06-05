package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/requests"
	"stockinos.com/api/services"
	"stockinos.com/api/utils"
)

type authInterface interface {
	CreateUserIfNotExists(ctx context.Context, phoneNumber string) error
	CreateOTP(ctx context.Context, pinCode models.OTP) error
	SaveOTP(ctx context.Context, pinCode models.OTP) error
	CheckOTP(ctx context.Context, phoneNumber, pinCode string) (*models.OTP, error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
}

type GetOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
}

type CheckOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
	PinCode     string `json:"pin_code,omitempty"`     // Pin code entered
}

func GetOTP(mux chi.Router, a authInterface) {
	mux.Post("/otp", func(w http.ResponseWriter, r *http.Request) {
		var input GetOTPRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number
		err := decoder.Decode(&input)
		log.Println("extract the phone number: ", err, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate the pin code of 4 digits
		now := time.Now()
		pinCode := utils.GenerateOTP(now)

		// send the pin code to a the phone number using Whatsapp API
		res, err := requests.SendWoZOTP(
			input.PhoneNumber,
			input.Language,
			pinCode,
		)
		if err != nil {
			log.Println("error when sending the OTP via WhatsApp: ", err)
			http.Error(w, "ERR_COTP_150", http.StatusBadRequest)
			return
		}

		// check if there is an user with this account
		err = a.CreateUserIfNotExists(r.Context(), input.PhoneNumber)
		if err != nil {
			log.Println("error when creating the user if he does not exist: ", err)
			http.Error(w, "ERR_COTP_151", http.StatusBadRequest)
			return
		}

		// if not, save the association phone number/pin code in the db
		var m models.OTP
		m.WaMessageId = res.Messages[0].ID
		m.PhoneNumber = input.PhoneNumber
		m.PinCode = pinCode
		m.Active = true

		err = a.CreateOTP(r.Context(), m)
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func CheckOTP(mux chi.Router, a authInterface) {
	mux.Post("/otp/check", func(w http.ResponseWriter, r *http.Request) {
		// read the request body
		var input CheckOTPRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number and the pin code
		err := decoder.Decode(&input)

		log.Println("extract the phone number: ", err, input)
		if err != nil {
			log.Println("error when extracting the request body: ", err)
			http.Error(w, "ERR_CTOP_101", http.StatusBadRequest)
			return
		}

		// check that the pin code is 6 digit
		var m *models.OTP

		// check that the phone number is correct
		m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
		if err != nil {
			log.Println("error when checking the otp: ", err)
			http.Error(w, "ERR_COTP_102", http.StatusBadRequest)
			return
		}

		m.Active = false
		err = a.SaveOTP(r.Context(), *m)
		if err != nil {
			log.Println("error when changing the active state of the current OTP line: ", err)
			http.Error(w, "ERR_COTP_103", http.StatusBadRequest)
			return
		}

		// Generating the JWT Token
		u, err := a.FindUserByPhoneNumber(r.Context(), input.PhoneNumber)
		if err != nil {
			log.Println("error when looking for user: ", err)
			http.Error(w, "ERR_COTP_104", http.StatusBadRequest)
			return
		}

		var signInResult requests.SignInResult
		signInResult.Name = u.Name
		signInResult.PhoneNumber = u.PhoneNumber

		tokenString, err := services.GenerateJWTToken(structs.Map(signInResult))
		if err != nil {
			log.Println("error when generating jwt token ", err)
			http.Error(w, "ERR_COTP_105", http.StatusBadRequest)
			return
		}

		signInResult.Token = tokenString

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(signInResult); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

func ResendOTP(mux chi.Router, a authInterface) {
	mux.Post("/otp/resend", func(w http.ResponseWriter, r *http.Request) {
		// // read the request body
		// var input CheckOTPRequest

		// // read the request body
		// decoder := json.NewDecoder(r.Body)

		// // extract the phone number and the pin code
		// err := decoder.Decode(&input)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// // check that the pin code is 6 digit
		// var m *models.OTP

		// // check that the phone number is correct
		// m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusBadRequest)
		// 	return
		// }

		// m.Active = false
		// a.SaveOTP(r.Context(), *m)
	})
}
