package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

type createOTPInterface interface {
	CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOTPx(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error)
}

type CreateOTPRequest struct {
	PhoneNumber string `json:"phone_number,omitempty"` // Phone number of the customer
	Language    string `json:"language,omitempty"`     // Language for template
}

type CreateOTPResponse struct {
	PhoneNumber string `json"phone_number"`
}

func (appHandler *AppHandler) CreateOTP(mux chi.Router, db createOTPInterface) {
	mux.Post("/otp", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input CreateOTPRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
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
		user, err := db.DoesUserExists(ctx, storage.DoesUserExistsParams{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			// log.Println("Error at DoesUserExists", err)
			http.Error(w, "ERR_AUTH_CRT_OTP_01", http.StatusBadRequest)
			return
		}

		// If there is none, we create the user
		if user == nil {
			_, err := db.CreateUser(ctx, storage.CreateUserParams{
				PhoneNumber: input.PhoneNumber,
				FirstName:   "",
				LastName:    "",
				Email:       "",
			})
			if err != nil {
				// log.Println("Error at CreateUser", err)
				http.Error(w, "ERR_AUTH_CRT_OTP_02", http.StatusBadRequest)
				return
			}
		}

		// Then we save the OTP
		// 1- update all otps to not active
		// 2- create the new otp as active
		_, err = db.CreateOTPx(ctx, storage.CreateOTPParams{
			WaMessageId: waMessageId,
			PhoneNumber: input.PhoneNumber,
			PinCode:     pinCode,
		})
		if err != nil {
			// log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_AUTH_CRT_OTP_03", http.StatusBadRequest)
			return
		}

		response := CreateOTPResponse{
			PhoneNumber: input.PhoneNumber,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_AUTH_CRT_OTP_END", http.StatusBadRequest)
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

type CheckOTPResponse struct {
	User models.User `json:"user"`
}

func (appHandler *AppHandler) CheckOTP(mux chi.Router, db checkOTPInterface) {
	mux.Post("/otp/check", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// read the request body
		var input CheckOTPRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		user, err := db.DoesUserExists(ctx, storage.DoesUserExistsParams{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_AUTH_CHK_OTP_01", http.StatusBadRequest)
			return
		}
		if user == nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_AUTH_CHK_OTP_02", http.StatusBadRequest)
			return
		}
		// check that the pin code is 6 digit
		// var m *models.OTP

		// check that the phone number is correct
		otp, err := db.CheckOTPTx(r.Context(), storage.CheckOTPParams{
			PhoneNumber: input.PhoneNumber,
			UserOTP:     input.PinCode,
		})
		if err != nil {
			// log.Println("error when checking the otp: ", err)
			http.Error(w, "ERR_AUTH_CHK_OTP_03", http.StatusBadRequest)
			return
		}
		if otp == nil {
			// log.Println("error when checking the otp - no corresponding otp found: ", err)
			http.Error(w, "ERR_AUTH_CHK_OTP_04", http.StatusBadRequest)
			return
		}

		tokenString, err := services.GenerateJwtToken(structs.Map(&user))
		if err != nil {
			log.Println("Error CreateUser", zap.Error(err))
			http.Error(w, "error when creating token", http.StatusBadRequest)
			return
		}

		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Secure:   true,
				MaxAge:   3600,
				SameSite: http.SameSiteLaxMode,
			},
		)

		response := CheckOTPResponse{
			User: *user,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_AUTH_CHK_OTP_END", http.StatusBadRequest)
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
