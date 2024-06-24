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
)

type getOTPMock struct{}

func TestOTP(t *testing.T) {
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
		"CreateOTP": testCreateOTP,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockCreateOTPDB struct {
	CreateUserFunc     func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error)
	DoesUserExistsFunc func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CreateOTPxFunc     func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error)
}

func (mdb *mockCreateOTPDB) CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
	return mdb.CreateUserFunc(ctx, arg)
}

func (mdb *mockCreateOTPDB) DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
	return mdb.DoesUserExistsFunc(ctx, arg)
}

func (mdb *mockCreateOTPDB) CreateOTPx(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
	return mdb.CreateOTPxFunc(ctx, arg)
}

func testCreateOTP(t *testing.T, appHandler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateOTPDB{
			CreateUserFunc: func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
				return nil, nil
			},
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return nil, nil
			},
			CreateOTPxFunc: func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CreateOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CreateOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CreateOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("error when checking if user exists", func(t *testing.T) {
		dataRequest := handlers.CreateOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
		}

		mux := chi.NewMux()
		db := &mockCreateOTPDB{
			CreateUserFunc: func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
				return nil, nil
			},
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return nil, errors.New("error when checking")
			},
			CreateOTPxFunc: func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CreateOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CreateOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CRT_OTP_01"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CreateOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("user does not exists", func(t *testing.T) {
		dataRequest := handlers.CreateOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
		}

		mux := chi.NewMux()

		t.Run("error when creating user", func(t *testing.T) {
			db := &mockCreateOTPDB{
				CreateUserFunc: func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
					return nil, errors.New("error when creating user")
				},
				DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
					return nil, nil
				},
				CreateOTPxFunc: func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
					return nil, nil
				},
			}

			appHandler.CreateOTP(mux, db)
			code, _, response := helpertest.MakePostRequest(
				mux,
				"/otp",
				helpertest.CreateFormHeader(),
				dataRequest,
				[]helpertest.ContextData{},
			)
			wantStatus := http.StatusBadRequest
			if code != wantStatus {
				t.Fatalf("CreateOTP(): status - got %d; want %d", code, wantStatus)
			}
			wantCode := "ERR_AUTH_CRT_OTP_02"
			if !strings.HasPrefix(response, wantCode) {
				t.Fatalf("CreateOTP(): response error - got %s, want %s", response, wantCode)
			}
		})
	})

	t.Run("error when creating OTP", func(t *testing.T) {
		dataRequest := handlers.CreateOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
		}

		mux := chi.NewMux()
		db := &mockCreateOTPDB{
			CreateUserFunc: func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
				return nil, nil
			},
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return &models.User{
					PhoneNumber: dataRequest.PhoneNumber,
					FirstName:   gofaker.FirstName(),
					LastName:    gofaker.LastName(),
					Email:       gofaker.Email(),
				}, nil
			},
			CreateOTPxFunc: func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
				return nil, errors.New("error when creating OTP")
			},
		}

		appHandler.CreateOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CreateOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CRT_OTP_03"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CreateOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.CreateOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
		}

		mux := chi.NewMux()
		db := &mockCreateOTPDB{
			CreateUserFunc: func(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
				return nil, nil
			},
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return &models.User{
					PhoneNumber: dataRequest.PhoneNumber,
					FirstName:   gofaker.FirstName(),
					LastName:    gofaker.LastName(),
					Email:       gofaker.Email(),
				}, nil
			},
			CreateOTPxFunc: func(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CreateOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusOK
		if code != wantStatus {
			t.Fatalf("CreateOTP(): status - got %d; want %d", code, wantStatus)
		}

		got := handlers.CreateOTPResponse{}
		json.Unmarshal([]byte(response), &got)
		if got.PhoneNumber != dataRequest.PhoneNumber {
			t.Fatalf("CreateOTP(): response Phone number - got %s; want %s", got.PhoneNumber, dataRequest.PhoneNumber)
		}
	})

}
