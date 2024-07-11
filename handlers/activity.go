package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type organizationMiddlewareInterface interface {
	GetActivity(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error)
}

func (handler *AppHandler) ActivityMiddleware(mux chi.Router, db organizationMiddlewareInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			activityIdParam := chi.URLParamFromCtx(ctx, "activityId")
			activityId, err := primitive.ObjectIDFromHex(activityIdParam)
			if err != nil {
				http.Error(w, "ERR_ATVT_MDW_01", http.StatusBadRequest)
				return
			}

			organization := ctx.Value("organization").(*models.Organization)

			activity, err := db.GetActivity(ctx, storage.GetActivityParams{
				Id:             activityId,
				OrganizationId: organization.Id,
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

		organization := ctx.Value("organization").(*models.Organization)

		activities, err := db.GetAllActivities(ctx, storage.GetAllActivitiesParams{
			OrganizationId: organization.Id,
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
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Fields      []models.ActivityField `json:"fields"`
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

		organization := ctx.Value("organization").(*models.Organization)

		var fields []models.ActivityField
		if len(input.Fields) == 0 {
			fields = make([]models.ActivityField, 0)
		} else {
			fields = input.Fields
		}

		activity, err := db.CreateActivity(ctx, storage.CreateActivityParams{
			Name:        input.Name,
			Description: input.Description,
			Fields:      fields,

			OrganizationId: organization.Id,
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

		organization := ctx.Value("organization").(*models.Organization)
		activity := ctx.Value("activity").(*models.Activity)

		err := db.DeleteActivity(ctx, storage.DeleteActivityParams{
			Id:             activity.Id,
			OrganizationId: organization.Id,
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
	UpdateSetInActivity(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error)
	UpdateAddToActivity(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error)
	UpdateRemoveFromActivity(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error)
}

type UpdateActivityRequest struct {
	Operation string      `json:"op"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
	Position  uint64      `json:"position"`
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

		organization := ctx.Value("organization").(*models.Organization)
		activity := ctx.Value("activity").(*models.Activity)

		var updatedActivity *models.Activity
		operation := strings.ToLower(input.Operation)
		field := strings.ToLower(input.Field)
		switch operation {
		case "set":
			if field == "" {
				http.Error(w, "ERR_ATVT_UDT_010", http.StatusBadRequest)
				return
			}
			var value any
			if field == "fields" {
				t, _ := json.Marshal(input.Value)
				j := []byte(t)
				var v []models.ActivityField
				err := json.Unmarshal(j, &v)
				if err != nil {
					http.Error(w, "ERR_ATVT_UDT_014", http.StatusBadRequest)
					return
				}

				value = v
			} else {
				value = input.Value
			}
			// TODO: check type of input.Value string|int|bool

			updatedActivity, err = db.UpdateSetInActivity(ctx, storage.UpdateSetInActivityParams{
				Id:             activity.Id,
				OrganizationId: organization.Id,

				Field: field,
				Value: value, // input.Value,
			})
		case "add":
			if field != "fields" {
				http.Error(w, "ERR_ATVT_UDT_011", http.StatusBadRequest)
				return
			}
			if err, ok := input.Value.(models.ActivityField); !ok {
				log.Printf("Error : %+v\n%+v\n\n%+v\n\n", err, input.Value, input)
				// http.Error(w, "ERR_ATVT_UDT_014", http.StatusBadRequest)
				// return

			}
			value := models.ActivityField{
				Id:          primitive.NewObjectID(),
				Name:        getOrDefault(input.Value.(map[string]any), "name", "").(string),
				Description: getOrDefault(input.Value.(map[string]any), "description", "").(string),
				Type:        getOrDefault(input.Value.(map[string]any), "type", "").(string),
				Key:         getOrDefault(input.Value.(map[string]any), "id", false).(bool),
				Code:        getOrDefault(input.Value.(map[string]any), "code", "").(string),
				Options: models.ActivityFieldOptions{
					Reference:    nil,
					DefaultValue: nil,
					Multiple:     false,
					Automatic:    false,
				},
			}

			updatedActivity, err = db.UpdateAddToActivity(ctx, storage.UpdateAddToActivityParams{
				Id:             activity.Id,
				OrganizationId: organization.Id,

				Field:    field,
				Value:    value,
				Position: uint(input.Position),
			})
		case "remove":
			// if match, _ := regexp.MatchString("fields.([0-9]+)", field); !match {
			if field != "fields" {
				http.Error(w, "ERR_ATVT_UDT_012", http.StatusBadRequest)
				return
			}
			// TODO: check type of input.Value Should be emtpy

			updatedActivity, err = db.UpdateRemoveFromActivity(ctx, storage.UpdateRemoveFromActivityParams{
				Id:             activity.Id,
				OrganizationId: organization.Id,

				Field:    field,
				Position: uint(input.Position),
			})
		default:
			http.Error(w, "ERR_ATVT_UDT_013", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "ERR_ATVT_UDT_02", http.StatusBadRequest)
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

func getOrDefault(m map[string]any, key string, defaultValue any) any {
	value, ok := m[key]
	if ok {
		return value
	} else {
		return defaultValue
	}
}
