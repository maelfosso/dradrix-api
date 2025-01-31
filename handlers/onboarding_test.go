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

func TestOnboarding(t *testing.T) {
	handler := handlers.NewAppHandler()
	handler.GetAuthenticatedUser = func(req *http.Request) *models.User {
		return authenticatedUser
	}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		"SetProfile":        testSetProfile,
		"FirstOrganization": testFirstOrganization,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockSetProfileDB struct {
	UpdateUserProfileFunc     func(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error)
	UpdateUserPreferencesFunc func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

func (mdb *mockSetProfileDB) UpdateUserProfile(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error) {
	return mdb.UpdateUserProfileFunc(ctx, arg)
}

func (mdb *mockSetProfileDB) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
	return mdb.UpdateUserPreferencesFunc(ctx, arg)
}

func testSetProfile(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockSetProfileDB{
			UpdateUserProfileFunc: func(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error) {
				return nil, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.SetProfile(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetProfile(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetProfile(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("error update user name", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockSetProfileDB{
			UpdateUserProfileFunc: func(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error) {
				return nil, errors.New("update user's name failed")
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.SetProfile(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			handlers.UpdateProfileRequest{
				FirstName: gofaker.FirstName(),
				LastName:  gofaker.LastName(),
			},
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetProfile(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_SN_01"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetProfile(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("error update user preferences", func(t *testing.T) {
		dataRequest := handlers.UpdateProfileRequest{
			FirstName: gofaker.FirstName(),
			LastName:  gofaker.LastName(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   dataRequest.FirstName,
			LastName:    dataRequest.LastName,

			Preferences: models.UserPreferences{
				Organization: models.UserPreferencesOrganization{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
				},
				OnboardingStep: 1,
			},
		}
		mux := chi.NewMux()
		db := &mockSetProfileDB{
			UpdateUserProfileFunc: func(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error) {
				return &updatedUser, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, errors.New("update's user preferences error")
			},
		}

		handler.SetProfile(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("SetProfile(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_SN_02"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("SetProfile(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.UpdateProfileRequest{
			FirstName: gofaker.FirstName(),
			LastName:  gofaker.LastName(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   dataRequest.FirstName,
			LastName:    dataRequest.LastName,

			Preferences: models.UserPreferences{
				Organization: models.UserPreferencesOrganization{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
				},
				OnboardingStep: 1,
			},
		}
		mux := chi.NewMux()
		db := &mockSetProfileDB{
			UpdateUserProfileFunc: func(ctx context.Context, arg storage.UpdateUserProfileParams) (*models.User, error) {
				return &updatedUser, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return &updatedUser, nil
			},
		}

		handler.SetProfile(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/name",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusOK
		if code != wantCode {
			t.Fatalf("SetProfile(): code - got %d; want %d", code, wantCode)
		}

		got := handlers.UpdateProfileResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Done {
			t.Fatalf("SetProfile(): response done - got %+v; want true", got.Done)
		}
	})
}

type mockFirstOrganizationDB struct {
	CreateOrganizationFunc    func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
	UpdateUserPreferencesFunc func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error)
}

func (mdb *mockFirstOrganizationDB) CreateOrganization(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
	return mdb.CreateOrganizationFunc(ctx, arg)
}

func (mdb *mockFirstOrganizationDB) UpdateUserPreferences(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
	return mdb.UpdateUserPreferencesFunc(ctx, arg)
}

func testFirstOrganization(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockFirstOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return nil, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.FirstOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/organization",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("FirstOrganization(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("FirstOrganization(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error creating organization", func(t *testing.T) {
		dataRequest := handlers.SetUpOrganizationRequest{
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		mux := chi.NewMux()
		db := &mockFirstOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return nil, errors.New("create organization error")
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, nil
			},
		}

		handler.FirstOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/organization",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("FirstOrganization(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_OBD_CPN_01"
		if response != want {
			t.Fatalf("FirstOrganization(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error update user preferences", func(t *testing.T) {
		dataRequest := handlers.SetUpOrganizationRequest{
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),

			CreatedBy: authenticatedUser.Id,
		}

		mux := chi.NewMux()
		db := &mockFirstOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return organization, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return nil, errors.New("update organization preferences error")
			},
		}

		handler.FirstOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/organization",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("FirstOrganization(): code - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_OBD_CPN_02"
		if !strings.HasPrefix(response, wantError) {
			t.Fatalf("FirstOrganization(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.SetUpOrganizationRequest{
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		organization := models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		updatedUser := models.User{
			Id:          authenticatedUser.Id,
			PhoneNumber: gofaker.Phonenumber(),
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),

			Preferences: models.UserPreferences{
				Organization: models.UserPreferencesOrganization{
					Id:   organization.Id,
					Name: organization.Bio,
				},
				OnboardingStep: 1,
			},
		}

		mux := chi.NewMux()
		db := &mockFirstOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return &organization, nil
			},
			UpdateUserPreferencesFunc: func(ctx context.Context, arg storage.UpdateUserPreferencesParams) (*models.User, error) {
				return &updatedUser, nil
			},
		}

		handler.FirstOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/organization",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("FirstOrganization(): status - got %d; want %d", code, want)
		}

		got := handlers.SetUpOrganizationResponse{}
		json.Unmarshal([]byte(response), &got)
		if got.Id != organization.Id {
			t.Fatalf("FirstOrganization(): response Id - got %s; want %s", got.Id, organization.Id)
		}
	})
}
