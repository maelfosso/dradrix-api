package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	gofaker "github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/handlers"
	"stockinos.com/api/helpertest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"

	sfaker "syreclabs.com/go/faker"
)

var authenticatedUser = &models.User{
	Id:          primitive.NewObjectID(),
	FirstName:   gofaker.FirstName(),
	LastName:    gofaker.LastName(),
	PhoneNumber: gofaker.Phonenumber(),
}

func TestOnboarding(t *testing.T) {
	handler := handlers.NewAppHandler()
	handler.GetAuthenticatedUser = func(req *http.Request) *models.User {
		return authenticatedUser
	}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		"SetName": testSetName,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockSetNameDB struct {
	UpdateUserNameFunc        func(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error)
	UpdateUserPreferencesFunc func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

func (mdb *mockSetNameDB) UpdateUserName(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error) {
	return mdb.UpdateUserNameFunc(ctx, arg)
}

func (mdb *mockSetNameDB) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
	return mdb.UpdateUserPreferencesFunc(ctx, arg)
}

func testSetName(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockSetNameDB{
			UpdateUserNameFunc: func(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error) {
				return nil, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.SetName(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetName(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetName(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("error update user name", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockSetNameDB{
			UpdateUserNameFunc: func(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error) {
				return nil, errors.New("update user's name failed")
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.SetName(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			handlers.SetNameRequest{
				FirstName: gofaker.FirstName(),
				LastName:  gofaker.LastName(),
			},
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetName(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_SN_01"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetName(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("error update user preferences", func(t *testing.T) {
		dataRequest := handlers.SetNameRequest{
			FirstName: gofaker.FirstName(),
			LastName:  gofaker.LastName(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   dataRequest.FirstName,
			LastName:    dataRequest.LastName,

			Preferences: models.UserPreferences{
				Company: models.UserPreferencesCompany{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
				},
				CurrentOnboardingStep: 1,
			},
		}
		mux := chi.NewMux()
		db := &mockSetNameDB{
			UpdateUserNameFunc: func(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error) {
				return &updatedUser, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, errors.New("update's user preferences error")
			},
		}

		handler.SetName(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetName(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_SN_02"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetName(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.SetNameRequest{
			FirstName: gofaker.FirstName(),
			LastName:  gofaker.LastName(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   dataRequest.FirstName,
			LastName:    dataRequest.LastName,

			Preferences: models.UserPreferences{
				Company: models.UserPreferencesCompany{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
				},
				CurrentOnboardingStep: 1,
			},
		}
		mux := chi.NewMux()
		db := &mockSetNameDB{
			UpdateUserNameFunc: func(ctx context.Context, arg storage.UpdateUserNameParams) (*models.User, error) {
				return &updatedUser, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return &updatedUser, nil
			},
		}

		handler.SetName(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusOK
		if code != wantCode {
			t.Fatalf("SetName(): code - got %d; want %d", code, wantCode)
		}

		got := handlers.SetNameResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Done {
			t.Fatalf("SetName(): response done - got %+v; want true", got.Done)
		}
	})
}
