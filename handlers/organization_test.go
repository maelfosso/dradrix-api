package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func TestOrganization(t *testing.T) {
	handler := handlers.NewAppHandler()
	handler.GetAuthenticatedUser = func(req *http.Request) *models.User {
		return &models.User{
			Id:          primitive.NewObjectID(),
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),
			PhoneNumber: gofaker.Phonenumber(),
		}
	}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		"OrganizationMiddleware": testOrganizationMiddleware,
		"GetAllCompanies":        testGetAllCompanies,
		"GetOrganization":        testGetOrganization,
		"CreateOrganization":     testCreateOrganization,
		"UpdateOrganization":     testUpdateOrganization,
		"DeleteOrganization":     testDeleteOrganization,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockOrganizationMiddlewareDB struct {
	GetOrganizationFunc func(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error)
}

func (mdb *mockOrganizationMiddlewareDB) GetOrganization(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error) {
	return mdb.GetOrganizationFunc(ctx, arg)
}

func testOrganizationMiddleware(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid organization id", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockOrganizationMiddlewareDB{
			GetOrganizationFunc: func(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error) {
				return nil, nil
			},
		}

		mux.Route("/{organizationId}", func(r chi.Router) {
			handler.OrganizationMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		})

		_, w, response := helpertest.MakeGetRequest(mux, "/1", []helpertest.ContextData{})
		code := w.StatusCode
		if code != http.StatusBadRequest {
			t.Fatalf("OrganizationMiddleware(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_CMP_MDW_01"
		if response != want {
			t.Fatalf("OrganizationMiddleware(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockOrganizationMiddlewareDB{
			GetOrganizationFunc: func(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error) {
				return nil, errors.New("an error happens")
			},
		}

		mux.Route("/{organizationId}", func(r chi.Router) {
			handler.OrganizationMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		})

		_, w, response := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{},
		)
		code := w.StatusCode
		if code != http.StatusBadRequest {
			t.Fatalf("OrganizationMiddleware(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_CMP_MDW_02"
		if response != want {
			t.Fatalf("OrganizationMiddleware(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("no organization found", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockOrganizationMiddlewareDB{
			GetOrganizationFunc: func(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error) {
				return nil, nil
			},
		}

		mux.Route("/{organizationId}", func(r chi.Router) {
			handler.OrganizationMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		})

		_, w, response := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{},
		)
		code := w.StatusCode

		if code != http.StatusNotFound {
			t.Fatalf("OrganizationMiddleware(): status - got %d; want %d", code, http.StatusNotFound)
		}
		want := "ERR_CMP_MDW_03"
		if response != want {
			t.Fatalf("OrganizationMiddleware(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("organization found", func(t *testing.T) {
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		mux := chi.NewMux()
		db := &mockOrganizationMiddlewareDB{
			GetOrganizationFunc: func(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error) {
				return organization, nil
			},
		}

		mux.Route("/{organizationId}", func(r chi.Router) {
			handler.OrganizationMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				got := r.Context().Value("organization").(*models.Organization)
				if err := organizationEq(got, organization); err != nil {
					t.Fatalf("OrganizationMiddleware(): %v", err)
				}
			})
		})

		_, w, _ := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{},
		)
		code := w.StatusCode
		if code != http.StatusOK {
			t.Fatalf("OrganizationMiddleware(): status - got %d; want %d", code, http.StatusOK)
		}
	})

}

type mockGetAllCompaniesDB struct {
	GetAllCompaniesFunc func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Organization, error)
}

func (mdb *mockGetAllCompaniesDB) GetAllCompanies(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Organization, error) {
	return mdb.GetAllCompaniesFunc(ctx, arg)
}

func testGetAllCompanies(t *testing.T, handler *handlers.AppHandler) {
	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockGetAllCompaniesDB{
			GetAllCompaniesFunc: func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Organization, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.GetAllCompanies(mux, db)
		_, w, response := helpertest.MakeGetRequest(mux, "/", []helpertest.ContextData{})
		code := w.StatusCode
		if code != http.StatusBadRequest {
			t.Fatalf("GetAllCompanies(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_GALL_CMP_01"
		if response != want {
			t.Fatalf("GetAllCompanies(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		var organizations []*models.Organization

		const NUM_COMPANIES_CREATED = 3
		mux := chi.NewMux()

		db := &mockGetAllCompaniesDB{
			GetAllCompaniesFunc: func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Organization, error) {
				for i := 0; i < NUM_COMPANIES_CREATED; i++ {
					organization := &models.Organization{
						Id:   primitive.NewObjectID(),
						Name: sfaker.Company().Name(),
						Bio:  gofaker.Paragraph(),
					}
					organizations = append(organizations, organization)
				}
				return organizations, nil
			},
		}

		handler.GetAllCompanies(mux, db)
		_, w, response := helpertest.MakeGetRequest(mux, "/", []helpertest.ContextData{})
		code := w.StatusCode
		if code != http.StatusOK {
			t.Fatalf("GetAllCompanies(): status - got %d; want %d", code, http.StatusOK)
		}

		got := handlers.GetAllCompaniesResponse{}
		json.Unmarshal([]byte(response), &got)
		for i, c := range got.Companies {
			if err := organizationEq(c, organizations[i]); err != nil {
				t.Fatalf("GetOrganization(): %d - %v", i, err)
			}
		}
	})
}

func testGetOrganization(t *testing.T, handler *handlers.AppHandler) {

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}

		db := &struct{}{}

		handler.GetOrganization(mux, db)
		_, w, response := helpertest.MakeGetRequest(
			mux,
			"/",
			[]helpertest.ContextData{
				{
					Name:  "organization",
					Value: organization,
				},
			},
		)
		code := w.StatusCode
		if code != http.StatusOK {
			t.Fatalf("GetOrganization(): status - got %d; want %d", code, http.StatusOK)
		}

		got := handlers.GetOrganizationResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := organizationEq(&got.Organization, organization); err != nil {
			t.Fatalf("GetOrganization(): %v", err)
		}
	})
}

type mockCreateOrganizationDB struct {
	CreateOrganizationFunc func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
}

func (mdb *mockCreateOrganizationDB) CreateOrganization(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
	return mdb.CreateOrganizationFunc(ctx, arg)
}

func testCreateOrganization(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return nil, nil
			},
		}

		handler.CreateOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateOrganization(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("CreateComponey(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.CreateOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateOrganizationRequest{
				Name: sfaker.Company().Name(),
				Bio:  gofaker.Paragraph(),
			},
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateOrganization(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_C_CMP_01"
		if response != want {
			t.Fatalf("CreateOrganization(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}

		mux := chi.NewMux()
		db := &mockCreateOrganizationDB{
			CreateOrganizationFunc: func(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error) {
				return organization, nil
			},
		}

		handler.CreateOrganization(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateOrganizationRequest{
				Name: organization.Name,
				Bio:  organization.Bio,
			},
			[]helpertest.ContextData{},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("CreateOrganization(): status - got %d; want %d", code, want)
		}

		got := handlers.CreateOrganizationResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := organizationEq(&got.Organization, organization); err != nil {
			t.Fatalf("GetOrganization(): %v", err)
		}
	})
}

type mockUpdateOrganizationDB struct {
	UpdateOrganizationFunc func(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error)
}

func (mdb *mockUpdateOrganizationDB) UpdateOrganization(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error) {
	return mdb.UpdateOrganizationFunc(ctx, arg)
}

func testUpdateOrganization(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateOrganizationDB{
			UpdateOrganizationFunc: func(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error) {
				return nil, nil
			},
		}

		handler.UpdateOrganization(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("UpdateOrganization(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("UpdateComponey(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateOrganizationDB{
			UpdateOrganizationFunc: func(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.UpdateOrganization(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.UpdateOrganizationRequest{
				Name: sfaker.Company().Name(),
				Bio:  gofaker.Paragraph(),
			},
			[]helpertest.ContextData{
				{Name: "organization", Value: &models.Organization{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
					Bio:  gofaker.Paragraph(),
				}},
			},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("UpdateOrganization(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_U_CMP_02"
		if response != wantError {
			t.Fatalf("UpdateOrganization(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}

		mux := chi.NewMux()
		db := &mockUpdateOrganizationDB{
			UpdateOrganizationFunc: func(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error) {
				return &models.Organization{
					Id:   organization.Id,
					Name: sfaker.Company().Name(),
					Bio:  gofaker.Paragraph(),
				}, nil
			},
		}

		handler.UpdateOrganization(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.UpdateOrganizationRequest{
				Name: organization.Name,
				Bio:  organization.Bio,
			},
			[]helpertest.ContextData{
				{
					Name: "organization", Value: organization,
				},
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("UpdateOrganization(): status - got %d; want %d", code, want)
		}

		got := handlers.UpdateOrganizationResponse{}
		json.Unmarshal([]byte(response), &got)
		if got.Organization.Id != organization.Id {
			t.Fatalf("UpdatedOrganization(): Id - got %s; want %s", got.Organization.Id, organization.Id)
		}
	})
}

type mockDeleteOrganizationDB struct {
	DeleteOrganizationFunc func(ctx context.Context, arg storage.DeleteOrganizationParams) error
}

func (mdb *mockDeleteOrganizationDB) DeleteOrganization(ctx context.Context, arg storage.DeleteOrganizationParams) error {
	return mdb.DeleteOrganizationFunc(ctx, arg)
}

func testDeleteOrganization(t *testing.T, handler *handlers.AppHandler) {
	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockDeleteOrganizationDB{
			DeleteOrganizationFunc: func(ctx context.Context, arg storage.DeleteOrganizationParams) error {
				return errors.New("an error happens")
			},
		}

		handler.DeleteOrganization(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			nil,
			[]helpertest.ContextData{
				{Name: "organization", Value: &models.Organization{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
					Bio:  gofaker.Paragraph(),
				}},
			},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("DeleteOrganization(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_D_CMP_02"
		if response != wantError {
			t.Fatalf("DeleteOrganization(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockDeleteOrganizationDB{
			DeleteOrganizationFunc: func(ctx context.Context, arg storage.DeleteOrganizationParams) error {
				return nil
			},
		}

		handler.DeleteOrganization(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			nil,
			[]helpertest.ContextData{
				{Name: "organization", Value: &models.Organization{
					Id:   primitive.NewObjectID(),
					Name: sfaker.Company().Name(),
					Bio:  gofaker.Paragraph(),
				}},
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("DeleteOrganization(): status - got %d; want %d", code, want)
		}

		got := handlers.DeleteOrganizationResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Deleted {
			t.Fatalf("DeleteOrganization(): got %v; want %v", got.Deleted, true)
		}
	})
}

func organizationEq(got, want *models.Organization) error {
	if got == want {
		return nil
	}
	if got == nil {
		return fmt.Errorf("got nil; want %v", want)
	}
	if want == nil {
		return fmt.Errorf("got %v; want nil", got)
	}
	if got.Id != want.Id {
		return fmt.Errorf("got.Id = %s; want %s", got.Id, want.Id)
	}
	if got.Name != want.Name {
		return fmt.Errorf("got.Name = %s; want %s", got.Name, want.Name)
	}
	if got.Bio != want.Bio {
		return fmt.Errorf("got.Bio = %s; want %s", got.Bio, want.Bio)
	}
	if got.CreatedBy != want.CreatedBy {
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
