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

type companyMiddlewareInterface interface {
	GetActivity(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error)
}

func (handler *AppHandler) ActivityMiddleware(mux chi.Router, db companyMiddlewareInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			activityIdParam := chi.URLParamFromCtx(ctx, "activityId")
			activityId, err := primitive.ObjectIDFromHex(activityIdParam)
			if err != nil {
				http.Error(w, "ERR_ATVT_MDW_01", http.StatusBadRequest)
				return
			}

			company := ctx.Value("company").(*models.Company)

			activity, err := db.GetActivity(ctx, storage.GetActivityParams{
				Id:        activityId,
				CompanyId: company.Id,
			})
			if err != nil {
				http.Error(w, "ERR_ATVT_MDW_02", http.StatusBadRequest)
				return
			}
			if activity == nil {
				http.Error(w, "ERR_ATVT_MDW_03", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "activity", activity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

type getAllActivitiesInterface interface {
	GetAllActivities(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error)
}

type GetAllActivitiesResponse struct {
	Activities []*models.Activity `json:"activities,omitempty"`
}

func (handler *AppHandler) GetAllActivities(mux chi.Router, db getAllActivitiesInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		company := ctx.Value("company").(*models.Company)

		activities, err := db.GetAllActivities(ctx, storage.GetAllActivitiesParams{
			CompanyId: company.Id,
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

type createActivityInterface interface {
	CreateActivity(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error)
}

type CreateActivityRequest struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Fields      []models.ActivityFields `json:"fields,omitempty"`
}

type CreateActivityResponse struct {
	Activity models.Activity `json:"activity"`
}

func (handler *AppHandler) CreateActivity(mux chi.Router, db createActivityInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input CreateActivityRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		company := ctx.Value("company").(*models.Company)

		activity, err := db.CreateActivity(ctx, storage.CreateActivityParams{
			Name:        input.Name,
			Description: input.Description,
			Fields:      input.Fields,

			CompanyId: company.Id,
		})
		if err != nil {
			http.Error(w, "ERR_ATVT_CRT_01", http.StatusBadRequest)
			return
		}

		response := CreateActivityResponse{
			Activity: *activity,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_ATVT_CRT_END", http.StatusBadRequest)
			return
		}
	})
}

type getActivityInterface interface {
}

type GetActivityResponse struct {
	Activity models.Activity
}

func (handler *AppHandler) GetActivity(mux chi.Router, db getActivityInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		activity := ctx.Value("activity").(*models.Activity)

		response := GetActivityResponse{
			Activity: *activity,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_ATVT_GONE_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteActivityInterface interface {
	DeleteActivity(ctx context.Context, arg storage.DeleteActivityParams) error
}

type DeleteActivityResponse struct {
	Deleted bool `json:"deleted"`
}

func (handler *AppHandler) DeleteActivity(mux chi.Router, db deleteActivityInterface) {
	mux.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		company := ctx.Value("company").(*models.Company)
		activity := ctx.Value("activity").(*models.Activity)

		err := db.DeleteActivity(ctx, storage.DeleteActivityParams{
			Id:        activity.Id,
			CompanyId: company.Id,
		})
		if err != nil {
			http.Error(w, "ERR_ATVT_DLT_01", http.StatusBadRequest)
			return
		}

		response := DeleteActivityResponse{
			Deleted: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_ATVT_DLT_END", http.StatusBadRequest)
			return
		}
	})
}

type updateActivityInterface interface {
	UpdateActivity(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error)
}

type UpdateActivityRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
type UpdateActivityResponse struct {
	Activity models.Activity `json:"activity"`
}

func (handler *AppHandler) UpdateActivity(mux chi.Router, db updateActivityInterface) {
	mux.Patch("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input UpdateActivityRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		company := ctx.Value("company").(*models.Company)
		activity := ctx.Value("activity").(*models.Activity)

		updatedActivity, err := db.UpdateActivity(ctx, storage.UpdateActivityParams{
			Id:        activity.Id,
			CompanyId: company.Id,

			Field: input.Field,
			Value: input.Value,
		})
		if err != nil {
			http.Error(w, "ERR_ATVT_UDT_01", http.StatusBadRequest)
			return
		}

		response := UpdateActivityResponse{
			Activity: *updatedActivity,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_ATVT_UDT_END", http.StatusBadRequest)
			return
		}
	})
}