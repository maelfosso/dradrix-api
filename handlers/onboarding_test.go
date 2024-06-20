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
		"SetName":         testSetName,
		"FirstCompany":    testFirstCompany,
		"EndOfOnboarding": testEndOfOnboarding,
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

type mockFirstCompanyDB struct {
	CreateCompanyFunc         func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error)
	UpdateUserPreferencesFunc func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

func (mdb *mockFirstCompanyDB) CreateCompany(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
	return mdb.CreateCompanyFunc(ctx, arg)
}

func (mdb *mockFirstCompanyDB) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
	return mdb.UpdateUserPreferencesFunc(ctx, arg)
}

func testFirstCompany(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockFirstCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return nil, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.FirstCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/company",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("FirstCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("FirstCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error creating company", func(t *testing.T) {
		dataRequest := handlers.FirstCompanyRequest{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		mux := chi.NewMux()
		db := &mockFirstCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return nil, errors.New("create company error")
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.FirstCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/company",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("FirstCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_OBD_CPN_01"
		if response != want {
			t.Fatalf("FirstCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error update user preferences", func(t *testing.T) {
		dataRequest := handlers.FirstCompanyRequest{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),

			CreatedBy: authenticatedUser.Id,
		}

		mux := chi.NewMux()
		db := &mockFirstCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return company, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, errors.New("update company preferences error")
			},
		}

		handler.FirstCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/company",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("FirstCompany(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_CPN_02"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("FirstCompany(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.FirstCompanyRequest{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		company := models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),

			Preferences: models.UserPreferences{
				Company: models.UserPreferencesCompany{
					Id:   company.Id,
					Name: company.Description,
				},
				CurrentOnboardingStep: 1,
			},
		}

		mux := chi.NewMux()
		db := &mockFirstCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return &company, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return &updatedUser, nil
			},
		}

		handler.FirstCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/company",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("FirstCompany(): status - got %d; want %d", code, want)
		}

		got := handlers.FirstCompanyResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Done {
			t.Fatalf("FirstCompany(): response done - got %+v; want true", got.Done)
		}
	})
}

type mockEndOfOnboardingDB struct {
	UpdateUserPreferencesFunc func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

func (mdb *mockEndOfOnboardingDB) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
	return mdb.UpdateUserPreferencesFunc(ctx, arg)
}

func testEndOfOnboarding(t *testing.T, handler *handlers.AppHandler) {
	t.Run("error update user preferences", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockEndOfOnboardingDB{
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, errors.New("update company preferences error")
			},
		}

		handler.EndOfOnboarding(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/end",
			helpertest.CreateFormHeader(),
			"{}",
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("EndOfOnboarding(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_END_01"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("EndOfOnboarding(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		company := models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),

			Preferences: models.UserPreferences{
				Company: models.UserPreferencesCompany{
					Id:   company.Id,
					Name: company.Description,
				},
				CurrentOnboardingStep: -1,
			},
		}

		mux := chi.NewMux()
		db := &mockEndOfOnboardingDB{
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return &updatedUser, nil
			},
		}

		handler.EndOfOnboarding(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/end",
			helpertest.CreateFormHeader(),
			"{}",
			[]helpertest.ContextData{},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("EndOfOnboarding(): status - got %d; want %d", code, want)
		}

		got := handlers.EndOfOnboardingResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Done {
			t.Fatalf("EndOfOnboarding(): response done - got %+v; want true", got.Done)
		}
	})
}
