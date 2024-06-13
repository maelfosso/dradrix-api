package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type getAllActivitiesInterface interface {
	GetAllActivities(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error)
}

type GetAllActivitiesResponse struct {
	Activities []*models.Activity `json:"activities,omitempty"`
}

func (handler *AppHandler) GetAllActivities(mux chi.Router, db getAllActivitiesInterface) {
	mux.Get("/activities", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		companyIdParam := chi.URLParamFromCtx(ctx, "companyId")
		companyId, err := primitive.ObjectIDFromHex(companyIdParam)
		if err != nil {
			http.Error(w, "ERR_GONE_CMP_01", http.StatusBadRequest)
			return
		}

		activities, err := db.GetAllActivities(ctx, storage.GetAllActivitiesParams{
			CompanyId: companyId,
		})
		if err != nil {
			http.Error(w, "ERR_ATVT_GALL_01", http.StatusBadRequest)
			return
		}

		response := GetAllActivitiesResponse{
			Activities: activities,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_ATVT_GALL_END", http.StatusBadRequest)
			return
		}
	})
}
