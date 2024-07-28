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

type dataMiddlewareInterface interface {
	GetData(ctx context.Context, arg storage.GetDataParams) (*models.Data, error)
}

func (handler *AppHandler) DataMiddleware(mux chi.Router, db dataMiddlewareInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			dataIdParam := chi.URLParamFromCtx(ctx, "dataId")
			dataId, err := primitive.ObjectIDFromHex(dataIdParam)
			if err != nil {
				http.Error(w, "ERR_DATA_MDW_01", http.StatusBadRequest)
				return
			}

			activity := ctx.Value("activity").(*models.Activity)

			data, err := db.GetData(ctx, storage.GetDataParams{
				Id:         dataId,
				ActivityId: activity.Id,
			})
			if err != nil {
				http.Error(w, "ERR_DATA_MDW_02", http.StatusBadRequest)
				return
			}
			if data == nil {
				http.Error(w, "ERR_DATA_MDW_03", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

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
		// We should ensure that all the data are the type of the one defined in activity
		// values := make(map[string]any)
		// for _, field := range activity.Fields {
		// 	value := input.Values[field.Code]
		// 	castValue, ok := field.IsValid(value)
		// 	if ok {
		// 		values[field.Code] = castValue
		// 	}
		// }
		// for code, value := range input.Values {
		// 	var field models.ActivityFields

		// 	for _, f := range activity.Fields {

		// 	}
		// }

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

type getAllDataInterface interface {
	GetAllData(ctx context.Context, arg storage.GetAllDataParams) ([]*models.Data, error)
}

type GetAllDataResponse struct {
	Fields map[string]string `json:"fields"`
	Data   []*models.Data    `json:"data"`
}

func (handler *AppHandler) GetAllData(mux chi.Router, db getAllDataInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		activity := ctx.Value("activity").(*models.Activity)

		data, err := db.GetAllData(ctx, storage.GetAllDataParams{
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_GALL_01", http.StatusBadRequest)
			return
		}

		fields := make(map[string]string)
		for _, field := range activity.Fields {
			fields[field.Id.Hex()] = field.Name
		}

		response := GetAllDataResponse{
			Fields: fields,
			Data:   data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_GALL_END", http.StatusBadRequest)
			return
		}
	})
}

type getDataInterface interface {
}

type GetDataResponse struct {
	Data models.Data
}

func (handler *AppHandler) GetData(mux chi.Router, db getDataInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := ctx.Value("data").(*models.Data)

		response := GetDataResponse{
			Data: *data,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_GONE_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteDataInterface interface {
	DeleteData(ctx context.Context, arg storage.DeleteDataParams) error
}

type DeleteDataResponse struct {
	Deleted bool `json:"deleted"`
}

func (handler *AppHandler) DeleteData(mux chi.Router, db deleteDataInterface) {
	mux.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		activity := ctx.Value("activity").(*models.Activity)
		data := ctx.Value("data").(*models.Data)

		err := db.DeleteData(ctx, storage.DeleteDataParams{
			Id:         data.Id,
			ActivityId: activity.Id,
		})
		if err != nil {
			http.Error(w, "ERR_DATA_DLT_01", http.StatusBadRequest)
			return
		}

		response := DeleteDataResponse{
			Deleted: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_DATA_DLT_END", http.StatusBadRequest)
			return
		}
	})
}
