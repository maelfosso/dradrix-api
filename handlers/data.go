package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type createDataInterface interface {
	CreateData(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error)
}

type CreateDataRequest struct {
	Values map[string]any `json:"values,omitempty"`
}

type CreateDataResponse struct {
	Data models.Data `json:"data"`
}

func (handler *AppHandler) CreateData(mux chi.Router, db createDataInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		var input CreateDataRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		activity := ctx.Value("activity").(*models.Activity)

		data, err := db.CreateData(ctx, storage.CreateDataParams{
			Values: input.Values,

			ActivityId: activity.Id,
			CreatedBy:  authUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_CRT_01", http.StatusBadRequest)
			return
		}

		response := CreateDataResponse{
			Data: *data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_CRT_END", http.StatusBadRequest)
			return
		}
	})
}
