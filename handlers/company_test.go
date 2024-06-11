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

func TestCompany(t *testing.T) {
	handler := handlers.NewAppHandler()
	handler.GetAuthenticatedUser = func(req *http.Request) *models.User {
		return &models.User{
			Id:          primitive.NewObjectID(),
			Name:        gofaker.Name(),
			PhoneNumber: gofaker.Phonenumber(),
		}
	}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		"GetAllCompanies": testGetAllCompanies,
		"GetCompany":      testGetCompany,
		"CreateCompany":   testCreateCompany,
		"UpdateCompany":   testUpdateCompany,
		"DeleteCompany":   testDeleteCompany,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockGetAllCompaniesDB struct {
	GetAllCompaniesFunc func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Company, error)
}

func (mdb *mockGetAllCompaniesDB) GetAllCompanies(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Company, error) {
	return mdb.GetAllCompaniesFunc(ctx, arg)
}

func testGetAllCompanies(t *testing.T, handler *handlers.AppHandler) {
	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockGetAllCompaniesDB{
			GetAllCompaniesFunc: func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Company, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.GetAllCompanies(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, "/")
		if code != http.StatusBadRequest {
			t.Fatalf("GetAllCompanies(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_GALL_CMP_01"
		if response != want {
			t.Fatalf("GetAllCompanies(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		var companies []*models.Company

		const NUM_COMPANIES_CREATED = 3
		mux := chi.NewMux()

		db := &mockGetAllCompaniesDB{
			GetAllCompaniesFunc: func(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Company, error) {
				for i := 0; i < NUM_COMPANIES_CREATED; i++ {
					company := &models.Company{
						Id:          primitive.NewObjectID(),
						Name:        sfaker.Company().Name(),
						Description: gofaker.Paragraph(),
					}
					companies = append(companies, company)
				}
				return companies, nil
			},
		}

		handler.GetAllCompanies(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, "/")
		if code != http.StatusOK {
			t.Fatalf("GetAllCompanies(): status - got %d; want %d", code, http.StatusOK)
		}

		got := handlers.GetAllCompaniesResponse{}
		json.Unmarshal([]byte(response), &got)
		for i, c := range got.Companies {
			if err := companyEq(c, companies[i]); err != nil {
				t.Fatalf("GetCompany(): %d - %v", i, err)
			}
		}
	})
}

type mockGetCompanyDB struct {
	GetCompanyFunc func(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error)
}

func (mdb *mockGetCompanyDB) GetCompany(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error) {
	return mdb.GetCompanyFunc(ctx, arg)
}

func testGetCompany(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid company id", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockGetCompanyDB{
			GetCompanyFunc: func(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error) {
				return nil, nil
			},
		}

		handler.GetCompany(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, "/1")
		if code != http.StatusBadRequest {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_GONE_CMP_01"
		if response != want {
			t.Fatalf("GetCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockGetCompanyDB{
			GetCompanyFunc: func(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.GetCompany(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, fmt.Sprintf("/%s", primitive.NewObjectID().Hex()))
		if code != http.StatusBadRequest {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_GONE_CMP_02"
		if response != want {
			t.Fatalf("GetCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("no company found", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockGetCompanyDB{
			GetCompanyFunc: func(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error) {
				return nil, nil
			},
		}

		handler.GetCompany(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, fmt.Sprintf("/%s", primitive.NewObjectID().Hex()))
		if code != http.StatusBadRequest {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_GONE_CMP_03"
		if response != want {
			t.Fatalf("GetCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}

		db := &mockGetCompanyDB{
			GetCompanyFunc: func(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error) {
				return company, nil
			},
		}

		handler.GetCompany(mux, db)
		code, _, response := helpertest.MakeGetRequest(mux, fmt.Sprintf("/%s", primitive.NewObjectID().Hex()))
		if code != http.StatusOK {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusOK)
		}

		got := handlers.GetCompanyResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := companyEq(&got.Company, company); err != nil {
			t.Fatalf("GetCompany(): %v", err)
		}
	})
}

type mockCreateCompanyDB struct {
	CreateCompanyFunc func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error)
}

func (mdb *mockCreateCompanyDB) CreateCompany(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
	return mdb.CreateCompanyFunc(ctx, arg)
}

func testCreateCompany(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return nil, nil
			},
		}

		handler.CreateCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("CreateComponey(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.CreateCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateCompanyRequest{
				Name:        sfaker.Company().Name(),
				Description: gofaker.Paragraph(),
			},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_C_CMP_01"
		if response != want {
			t.Fatalf("CreateCompany(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}

		mux := chi.NewMux()
		db := &mockCreateCompanyDB{
			CreateCompanyFunc: func(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error) {
				return company, nil
			},
		}

		handler.CreateCompany(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateCompanyRequest{
				Name:        company.Name,
				Description: company.Description,
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("CreateCompany(): status - got %d; want %d", code, want)
		}

		got := handlers.CreateCompanyResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := companyEq(&got.Company, company); err != nil {
			t.Fatalf("GetCompany(): %v", err)
		}
	})
}

type mockUpdateCompanyDB struct {
	UpdateCompanyFunc func(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error)
}

func (mdb *mockUpdateCompanyDB) UpdateCompany(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error) {
	return mdb.UpdateCompanyFunc(ctx, arg)
}

func testUpdateCompany(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid company id", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateCompanyDB{
			UpdateCompanyFunc: func(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error) {
				return nil, nil
			},
		}

		handler.UpdateCompany(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			"/1",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
		)
		if code != http.StatusBadRequest {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		wantCode := "ERR_U_CMP_01"
		if response != wantCode {
			t.Fatalf("GetCompany(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateCompanyDB{
			UpdateCompanyFunc: func(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error) {
				return nil, nil
			},
		}

		handler.UpdateCompany(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
		)
		if code != http.StatusBadRequest {
			t.Fatalf("UpdateCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("CreateComponey(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateCompanyDB{
			UpdateCompanyFunc: func(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.UpdateCompany(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			helpertest.CreateFormHeader(),
			handlers.UpdateCompanyRequest{
				Name:        sfaker.Company().Name(),
				Description: gofaker.Paragraph(),
			},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("UpdateCompany(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_U_CMP_02"
		if response != wantError {
			t.Fatalf("UpdateCompany(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}

		mux := chi.NewMux()
		db := &mockUpdateCompanyDB{
			UpdateCompanyFunc: func(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error) {
				return company, nil
			},
		}

		handler.UpdateCompany(mux, db)
		code, _, response := helpertest.MakePutRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			helpertest.CreateFormHeader(),
			handlers.UpdateCompanyRequest{
				Name:        company.Name,
				Description: company.Description,
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("UpdateCompany(): status - got %d; want %d", code, want)
		}

		got := handlers.UpdateCompanyResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := companyEq(&got.Company, company); err != nil {
			t.Fatalf("GetCompany(): %v", err)
		}
	})
}

type mockDeleteCompanyDB struct {
	DeleteCompanyFunc func(ctx context.Context, arg storage.DeleteCompanyParams) error
}

func (mdb *mockDeleteCompanyDB) DeleteCompany(ctx context.Context, arg storage.DeleteCompanyParams) error {
	return mdb.DeleteCompanyFunc(ctx, arg)
}

func testDeleteCompany(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid company id", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockDeleteCompanyDB{
			DeleteCompanyFunc: func(ctx context.Context, arg storage.DeleteCompanyParams) error {
				return nil
			},
		}

		handler.DeleteCompany(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			"/1",
			helpertest.CreateFormHeader(),
			nil,
		)
		if code != http.StatusBadRequest {
			t.Fatalf("GetCompany(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		wantCode := "ERR_D_CMP_01"
		if response != wantCode {
			t.Fatalf("GetCompany(): response error - got %s, want %s", response, wantCode)
		}
	})

	// t.Run("invalid input data", func(t *testing.T) {
	// 	mux := chi.NewMux()
	// 	db := &mockDeleteCompanyDB{
	// 		DeleteCompanyFunc: func(ctx context.Context, arg storage.DeleteCompanyParams) error {
	// 			return nil
	// 		},
	// 	}

	// 	handler.DeleteCompany(mux, db)
	// 	code, _, response := helpertest.MakeDeleteRequest(
	// 		mux,
	// 		fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
	// 		helpertest.CreateFormHeader(),
	// 		"{\"test\": \"that\"}",
	// 	)
	// 	if code != http.StatusBadRequest {
	// 		t.Fatalf("DeleteCompany(): status - got %d; want %d", code, http.StatusBadRequest)
	// 	}
	// 	want := "ERR_HDL_PRB_"
	// 	if !strings.HasPrefix(response, want) {
	// 		t.Fatalf("CreateComponey(): response error - got %s, want %s", response, want)
	// 	}
	// })

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockDeleteCompanyDB{
			DeleteCompanyFunc: func(ctx context.Context, arg storage.DeleteCompanyParams) error {
				return errors.New("an error happens")
			},
		}

		handler.DeleteCompany(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			helpertest.CreateFormHeader(),
			nil,
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("DeleteCompany(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_D_CMP_02"
		if response != wantError {
			t.Fatalf("DeleteCompany(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockDeleteCompanyDB{
			DeleteCompanyFunc: func(ctx context.Context, arg storage.DeleteCompanyParams) error {
				return nil
			},
		}

		handler.DeleteCompany(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			helpertest.CreateFormHeader(),
			nil,
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("DeleteCompany(): status - got %d; want %d", code, want)
		}

		got := handlers.DeleteCompanyResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Deleted {
			t.Fatalf("DeleteCompany(): got %v; want %v", got.Deleted, true)
		}
	})
}

func companyEq(got, want *models.Company) error {
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
	if got.Description != want.Description {
		return fmt.Errorf("got.Description = %s; want %s", got.Description, want.Description)
	}
	if got.CreatedBy != want.CreatedBy {
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
