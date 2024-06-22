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

type FirstCompanyInterface interface {
	CreateCompany(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error)
	UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

type FirstCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FirstCompanyResponse struct {
	Done bool `json:"done,omitempty"`
}

func (appHandler *AppHandler) FirstCompany(mux chi.Router, db FirstCompanyInterface) {
	mux.Post("/organization", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		currentAuthUser := appHandler.GetAuthenticatedUser(r)

		var input CreateCompanyRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		company, err := db.CreateCompany(ctx, storage.CreateCompanyParams{
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   currentAuthUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_01", http.StatusBadRequest)
			return
		}

		_, err = db.UpdateUserPreferences(ctx, storage.UpdateUserPreferencesParams{
			Id: currentAuthUser.Id,

			Changes: map[string]any{
				"company": map[string]any{
					"_id":  company.Id,
					"name": company.Name,
				},
				"onboarding_step": 3,
			},
		})
		if err != nil {
			http.Error(w, "ERR_OBD_CPN_02", http.StatusBadRequest)
			return
		}

		response := FirstCompanyResponse{
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
