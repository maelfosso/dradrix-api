package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type SetProfileInterface interface {
	UpdateUserProfile(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type SetProfileRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

type SetProfileResponse struct {
	Done bool `json:"done,omitempty"`
}

func (appHandler *AppHandler) SetProfile(mux chi.Router, db SetProfileInterface) {
	mux.Post("/profile", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input SetProfileRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		currentAuthUser := appHandler.GetAuthenticatedUser(r)

		_, err = db.UpdateUserProfile(ctx, storage.UpdateUserProfileParams{
			Id:        currentAuthUser.Id,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_SN_01", http.StatusBadRequest)
			return
		}

		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: currentAuthUser.Id,

			Changes: map[string]any{
				"onboarding_step": 1,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_SN_02", http.StatusBadRequest)
			return
		}

		response := SetProfileResponse{
			Done: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_OBD_SN_END", http.StatusBadRequest)
			return
		}
	})
}

type FirstOrganizationInterface interface {
	CreateOrganization(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type FirstOrganizationRequest struct {
	Name    string         `json:"name"`
	Bio     string         `json:"bio"`
	Email   string         `json:"email"`
	Address models.Address `json:"address"`
}

type FirstOrganizationResponse struct {
	Done bool `json:"done,omitempty"`
}

func (appHandler *AppHandler) FirstOrganization(mux chi.Router, db FirstOrganizationInterface) {
	mux.Post("/organization", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		currentAuthUser := appHandler.GetAuthenticatedUser(r)

		var input FirstOrganizationRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		organization, err := db.CreateOrganization(ctx, storage.CreateOrganizationParams{
			Name:      input.Name,
			Bio:       input.Bio,
			Email:     input.Email,
			Address:   input.Address,
			CreatedBy: currentAuthUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_01", http.StatusBadRequest)
			return
		}

		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: currentAuthUser.Id,

			Changes: map[string]any{
				"organization": map[string]any{
					"_id":  organization.Id,
					"name": organization.Name,
				},
				"onboarding_step": 2,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_02", http.StatusBadRequest)
			return
		}

		response := FirstOrganizationResponse{
			Done: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_OBD_CPN_END", http.StatusBadRequest)
			return
		}
	})
}

type EndOfOnboardingInterface interface {
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type EndOfOnboardingResponse struct {
	Done bool `json:"done,omitempty"`
}

func (appHandler *AppHandler) EndOfOnboarding(mux chi.Router, db EndOfOnboardingInterface) {
	mux.Post("/end", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		currentAuthUser := appHandler.GetAuthenticatedUser(r)

		_, err := db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: currentAuthUser.Id,
			Changes: map[string]any{
				"onboarding_step": -1,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_END_01", http.StatusBadRequest)
			return
		}

		response := EndOfOnboardingResponse{
			Done: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_OBD_END_END", http.StatusBadRequest)
			return
		}
	})
}
