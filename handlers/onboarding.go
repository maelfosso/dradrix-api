package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
)

type SetNameInterface interface {
	UpdateUserName(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error)
}

type SetNameRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type SetNameResponse struct {
	Saved bool `json:"saved,omitempty"`
}

func (appHandler *AppHandler) SetName(mux chi.Router, db SetNameInterface) {
	mux.Post("/name", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input SetNameRequest
		httpStatus, err := appHandler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		currentAuthUser := ctx.Value(services.JwtUserKey).(*models.User)

		_, err = db.UpdateUserName(ctx, storage.UpdateUserNameParams{
			Id:        currentAuthUser.Id,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		})
		if err != nil {
			http.Error(w, "ERR_OBD_SN_01", http.StatusBadRequest)
			return
		}

		response := SetNameResponse{
			Saved: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_GALL_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type FirstCompanyInterface interface{}

func (appHandler *AppHandler) FirstCompany(mux chi.Router, db FirstCompanyInterface) {
	mux.Post("/company", func(w http.ResponseWriter, r *http.Request) {

	})
}

type EndOfOnboardingInterface interface{}

func (appHandler *AppHandler) EndOfOnboarding(mux chi.Router, db EndOfOnboardingInterface) {
	mux.Post("/end", func(w http.ResponseWriter, r *http.Request) {

	})
}
