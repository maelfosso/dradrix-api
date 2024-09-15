package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

type getOTPInterface interface {
	CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOTPx(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type CreateOTPRequest struct {
	PhoneNumber string `json:"phone_number"` // Phone number of the customer
	Language    string `json:"language"`     // Language for template
}

type CreateOTPResponse struct {
	PhoneNumber   string `json:"phone_number"`
	RedirectToUrl string `json:"redirect_to_url"`
}

func (appHandler *AppHandler) CreateOTP(mux chi.Router, db getOTPInterface) {
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

		var user *models.User
		// Check if there is an user with this phone number
		user, err = db.DoesUserExists(ctx, storage.DoesUserExistsParams{
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			log.Println("Error at DoesUserExists", err)
			return
		}

		// If there is none, we create the user
		if user == nil {
			_, err = db.CreateUser(ctx, storage.CreateUserParams{
				PhoneNumber: input.PhoneNumber,
				FirstName:   "",
				LastName:    "",
			})
			if err != nil {
				log.Println("Error at CreateUser", err)
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
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		response := CreateOTPResponse{
			RedirectToUrl: fmt.Sprintf("/auth/check-otp?phone-number=%s", input.PhoneNumber),
			PhoneNumber:   input.PhoneNumber,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

type checkOTPInterface interface {
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CheckOTPTx(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type CheckOTPRequest struct {
	PhoneNumber string `json:"phone_number"` // Phone number of the customer
	Language    string `json:"language"`     // Language for template
	PinCode     string `json:"pin_code"`     // Pin code entered
}

type CheckOTPResponse struct {
	User          models.User `json:"user"`
	RedirectToUrl string      `json:"redirect_to_url"`
}

func (appHandler *AppHandler) CheckOTP(mux chi.Router, db checkOTPInterface) {
	mux.Post("/otp-check", func(w http.ResponseWriter, r *http.Request) {
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
		otp, err := db.CheckOTPTx(r.Context(), storage.CheckOTPParams{
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

		finalCurrentStatus := updateCurrentStatus(user.Preferences.CurrentStatus, "account-checked")
		user, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: user.Id,
			Changes: map[string]any{
				"current_status": finalCurrentStatus,
			},
		})
		if err != nil {
			http.Error(w, "ERR_COTP_ADD_MBR_PFRS", http.StatusBadRequest)
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
				Expires:  time.Now().Add(3600 * 24 * time.Second),
				MaxAge:   3600 * 24,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			},
		)

		reqOrigin := r.URL.Query().Get("from")
		redirectToUrl := nextLocation(reqOrigin, finalCurrentStatus)
		if redirectToUrl == "" {
			redirectToUrl = "/x" // fmt.Sprintf("/org/%s", user.Preferences.CurrentOrganizationId.Hex())
		} else {
			redirectToUrl = fmt.Sprintf("%s?phone-number=%s", redirectToUrl, input.PhoneNumber)
		}

		response := CheckOTPResponse{
			// User:          *user,
			RedirectToUrl: redirectToUrl,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("error when encoding auth result: ", err)
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}

type UpdateProfileInterface interface {
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	UpdateUserProfile(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type UpdateProfileRequest struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

type UpdateProfileResponse struct {
	Done          bool   `json:"done"`
	RedirectToUrl string `json:"redirect_to_url"`
}

func (appHandler *AppHandler) UpdateProfile(mux chi.Router, db UpdateProfileInterface) {
	mux.Post("/profile", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input UpdateProfileRequest
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
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}
		if user == nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		user, err = db.UpdateUserProfile(ctx, storage.UpdateUserProfileParams{
			Id:        user.Id,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_SN_01", http.StatusBadRequest)
			return
		}

		finalCurrentStatus := strings.Split(user.Preferences.CurrentStatus, "/")[0]
		if user.Preferences.CurrentOrganizationId.IsZero() {
			finalCurrentStatus = fmt.Sprintf("%s/set-org", finalCurrentStatus)
		} else {
			finalCurrentStatus = "registration-complete"
		}
		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: user.Id,

			Changes: map[string]any{
				"current_status": finalCurrentStatus,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_SN_02", http.StatusBadRequest)
			return
		}

		reqOrigin := r.URL.Query().Get("from")
		redirectToUrl := nextLocation(reqOrigin, finalCurrentStatus)
		if redirectToUrl == "" {
			redirectToUrl = "/x" // fmt.Sprintf("/org/%s", user.Preferences.CurrentOrganizationId.Hex())
		} else {
			redirectToUrl = fmt.Sprintf("%s?phone-number=%s", redirectToUrl, input.PhoneNumber)
		}
		response := UpdateProfileResponse{
			Done:          true,
			RedirectToUrl: redirectToUrl,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_OBD_SN_END", http.StatusBadRequest)
			return
		}
	})
}

type SetUpOrganizationInterface interface {
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOrganization(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
	AddMemberIntoOrganization(ctx context.Context, arg storage.AddMemberIntoOrganizationParams) (*models.Member, error)
}

type SetUpOrganizationRequest struct {
	PhoneNumber string         `json:"phone_number"`
	Name        string         `json:"name"`
	Bio         string         `json:"bio"`
	Email       string         `json:"email"`
	Address     models.Address `json:"address"`
}

type SetUpOrganizationResponse struct {
	Id            primitive.ObjectID `json:"id"`
	RedirectToUrl string             `json:"redirect_to_url"`
}

func (appHandler *AppHandler) SetUpOrganization(mux chi.Router, db SetUpOrganizationInterface) {
	mux.Post("/organization", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input SetUpOrganizationRequest
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
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}
		if user == nil {
			log.Println("error when saving the OTP: ", err)
			http.Error(w, "ERR_COTP_152", http.StatusBadRequest)
			return
		}

		organization, err := db.CreateOrganization(ctx, storage.CreateOrganizationParams{
			Name:            input.Name,
			Bio:             input.Bio,
			Email:           input.Email,
			Address:         input.Address,
			CreatedBy:       user.Id,
			OwnedBy:         user.Id,
			InvitationToken: uuid.New().String(),
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_01", http.StatusBadRequest)
			return
		}

		_, err = db.AddMemberIntoOrganization(ctx, storage.AddMemberIntoOrganizationParams{
			OrganizationId: organization.Id,
			UserId:         user.Id,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_12", http.StatusBadRequest)
			return
		}

		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: user.Id,

			Changes: map[string]any{
				"current_organization_id": organization.Id,
				"current_status":          "registration-complete",
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_02", http.StatusBadRequest)
			return
		}

		response := SetUpOrganizationResponse{
			Id:            organization.Id,
			RedirectToUrl: fmt.Sprintf("/org/%s", organization.Id.Hex()),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_OBD_CPN_END", http.StatusBadRequest)
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
