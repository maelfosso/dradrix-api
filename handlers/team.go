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

type getTeamInterface interface {
	GetMembersFromOrganization(ctx context.Context, arg storage.GetMembersFromOrganizationParams) ([]models.Member, error)
}

type GetTeamResponse struct {
	Members []models.Member
}

func (appHandler *AppHandler) GetTeam(mux chi.Router, db getTeamInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organization := ctx.Value("organization").(*models.Organization)

		members, err := db.GetMembersFromOrganization(ctx, storage.GetMembersFromOrganizationParams{
			OrganizationId: organization.Id,
		})
		if err != nil {
			http.Error(w, "ERR_TEAM_GTM_DB", http.StatusBadRequest)
			return
		}

		response := GetTeamResponse{
			Members: members,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_TEAM_GTM_END", http.StatusBadRequest)
			return
		}
	})
}

type invitationMiddlewareInterface interface {
	GetOrganizationFromInvitationToken(ctx context.Context, arg storage.GetOrganizationFromInvitationTokenParams) (*models.Organization, error)
}

func (handler *AppHandler) InvitationMiddleware(mux chi.Router, db invitationMiddlewareInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			invitationToken := chi.URLParamFromCtx(ctx, "invitationToken")

			organization, err := db.GetOrganizationFromInvitationToken(ctx, storage.GetOrganizationFromInvitationTokenParams{
				InvitationToken: invitationToken,
			})
			if err != nil {
				http.Error(w, "ERR_CMP_MDW_02", http.StatusBadRequest)
				return
			}

			if organization == nil {
				http.Error(w, "ERR_CMP_MDW_03", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "organization", organization)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

type getOrganizationFromInvitationTokenInterface interface{}

type GetOrganizationFromInvitationTokenResponse struct {
	Organization models.Organization `json:"organization,omitempty"`
}

func (handler *AppHandler) GetOrganizationFromInvitationToken(mux chi.Router, db getOrganizationFromInvitationTokenInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organization := ctx.Value("organization").(*models.Organization)

		response := GetOrganizationFromInvitationTokenResponse{
			Organization: *organization,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_C_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type addMemberInterface interface {
	CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOTPx(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error)
	AddMemberIntoOrganization(ctx context.Context, arg storage.AddMemberIntoOrganizationParams) (*models.Member, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type AddMemberRequest struct {
	PhoneNumber string `json:"phone_number"` // Phone number of the customer
	Language    string `json:"language"`     // Language for template
}

type AddMemberResponse struct {
	PhoneNumber   string `json:"phone_number"`
	RedirectToUrl string `json:"redirect_to_url"`
}

func (appHandler *AppHandler) AddMember(mux chi.Router, db addMemberInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		organization := ctx.Value("organization").(*models.Organization)

		var input AddMemberRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		// generate the pin code of 4 digits
		now := time.Now()
		pinCode := utils.GenerateOTP(now)

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
			user, err = db.CreateUser(ctx, storage.CreateUserParams{
				PhoneNumber: input.PhoneNumber,
				FirstName:   "",
				LastName:    "",
			})
			if err != nil {
				log.Println("Error at CreateUser", err)
				return
			}
		}

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

		_, err = db.AddMemberIntoOrganization(ctx, storage.AddMemberIntoOrganizationParams{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			InvitedAt:      now,
			ConfirmedAt:    &now,
		})
		if err != nil {
			http.Error(w, "ERR_COTP_ADD_MBR_ORG", http.StatusBadRequest)
			return
		}

		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: user.Id,

			Changes: map[string]any{
				"current_organization_id": organization.Id,
				"current_status":          user.Preferences.CurrentStatus,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_02", http.StatusBadRequest)
			return
		}

		redirectToUrl := fmt.Sprintf(
			"/join/%s/check-otp?phone-number=%s",
			organization.InvitationToken, input.PhoneNumber,
		)

		response := AddMemberResponse{
			RedirectToUrl: redirectToUrl,
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
