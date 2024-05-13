package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

type getOTPInterface interface {
	CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOTP(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error)
}

type CreateOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
}

func CreateOTP(mux chi.Router, svc getOTPInterface) {
	mux.Post("/otp", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input CreateOTPRequest

		// read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// generate the pin code of 4 digits
		now := time.Now()
		pinCode := utils.GenerateOTP(now)

		// send the pin code to a the phone number using Whatsapp API
		// waMessageId, err := requests.SendWoZOTP(
		// 	input.PhoneNumber,
		// 	input.Language,
		// 	pinCode,
		// )
		// if err != nil {
		// 	log.Println("error when sending the OTP via WhatsApp: ", err)
		// 	http.Error(w, "ERR_COTP_150", http.StatusBadRequest)
		// 	return
		// }
		waMessageId := "xxx-yyy-zzz"

		// Check if there is an user with this phone number
		user, err := svc.DoesUserExists(ctx, storage.DoesUserExistsParams{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			log.Println("Error at DoesUserExists", err)
			return
		}

		// If there is none, we create the user
		if user == nil {
			_, err := svc.CreateUser(ctx, storage.CreateUserParams{
				PhoneNumber: input.PhoneNumber,
				Name:        "",
			})
			if err != nil {
				log.Println("Error at CreateUser", err)
				return
			}
		}

		// Then we save the OTP
		_, err = svc.CreateOTP(ctx, storage.CreateOTPParams{
			WaMessageId: waMessageId,
			PhoneNumber: input.PhoneNumber,
			PinCode:     pinCode,
		})
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(input.PhoneNumber); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

type checkOTPInterface interface {
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CheckOTPTx(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error)
}

type CheckOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
	PinCode     string `json:"pin_code,omitempty"`     // Pin code entered
}

func CheckOTP(mux chi.Router, svc checkOTPInterface) {
	mux.Post("/otp/check", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// read the request body
		var input CheckOTPRequest

		// // read the request body
		decoder := json.NewDecoder(r.Body)

		// extract the phone number and the pin code
		err := decoder.Decode(&input)
		if err != nil {
			log.Println("error when extracting the request body: ", err)
			http.Error(w, "ERR_CTOP_101", http.StatusBadRequest)
			return
		}

		user, err := svc.DoesUserExists(ctx, storage.DoesUserExistsParams{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}
		if user == nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}
		// check that the pin code is 6 digit
		// var m *models.OTP

		// check that the phone number is correct
		otp, err := svc.CheckOTPTx(r.Context(), storage.CheckOTPParams{
			PhoneNumber: input.PhoneNumber,
			UserOTP:     input.PinCode,
		})
		if err != nil {
			log.Println("error when checking the otp: ", err)
			http.Error(w, fmt.Sprintf("ERR_COTP_102_%s", err), http.StatusBadRequest)
			return
		}
		if otp == nil {
			log.Println("error when checking the otp - no corresponding otp found: ", err)
			http.Error(w, "ERR_CHECK_OTP_", http.StatusBadRequest)
			return
		}

		// m.Active = false
		// err = a.SaveOTP(r.Context(), *m)
		// if err != nil {
		// 	log.Println("error when changing the active state of the current OTP line: ", err)
		// 	http.Error(w, "ERR_COTP_103", http.StatusBadRequest)
		// 	return
		// }

		// // Generating the JWT Token
		// u, err := a.FindUserByPhoneNumber(r.Context(), input.PhoneNumber)
		// if err != nil {
		// 	log.Println("error when looking for user: ", err)
		// 	http.Error(w, "ERR_COTP_104", http.StatusBadRequest)
		// 	return
		// }

		// var signInResult requests.SignInResult
		// signInResult.Name = u.Name
		// signInResult.PhoneNumber = u.PhoneNumber

		// tokenString, err := services.GenerateJWTToken(structs.Map(signInResult))
		// if err != nil {
		// 	log.Println("error when generating jwt token ", err)
		// 	http.Error(w, "ERR_COTP_105", http.StatusBadRequest)
		// 	return
		// }

		// signInResult.Token = tokenString

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(true); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

// func ResendOTP(mux chi.Router, a authInterface) {
// 	mux.Post("/otp/resend", func(w http.ResponseWriter, r *http.Request) {
// 		// // read the request body
// 		// var input CheckOTPRequest

// 		// // read the request body
// 		// decoder := json.NewDecoder(r.Body)

// 		// // extract the phone number and the pin code
// 		// err := decoder.Decode(&input)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 		// 	return
// 		// }

// 		// // check that the pin code is 6 digit
// 		// var m *models.OTP

// 		// // check that the phone number is correct
// 		// m, err = a.CheckOTP(r.Context(), input.PhoneNumber, input.PinCode)
// 		// if err != nil {
// 		// 	http.Error(w, err.Error(), http.StatusBadRequest)
// 		// 	return
// 		// }

// 		// m.Active = false
// 		// a.SaveOTP(r.Context(), *m)
// 	})
// }
