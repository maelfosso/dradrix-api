package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/handlers"
	"stockinos.com/api/helpertest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

func TestActivity(t *testing.T) {
	handler := &handlers.AppHandler{}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		"GetAllActivities": testGetAllActivities,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockGetAllActivities struct {
	GetAllActivitiesFunc func(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error)
}

func (mdb *mockGetAllActivities) GetAllActivities(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error) {
	return mdb.GetAllActivitiesFunc(ctx, arg)
}

func testGetAllActivities(t *testing.T, handler *handlers.AppHandler) {
	mockDb := &mockGetAllActivities{}

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()

		mockDb.GetAllActivitiesFunc = func(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error) {
			return []*models.Activity{}, nil
		}

		handler.GetAllActivities(mux, mockDb)
		_, w, response := helpertest.MakeGetRequest(mux, "/activities", []helpertest.ContextData{})
		code := w.StatusCode
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("GetAllActivities(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_ATVT_GALL_01"
		if response != wantError {
			t.Fatalf("GetAllActivities(): status - got %s; want %s", response, wantError)
		}
	})
}
