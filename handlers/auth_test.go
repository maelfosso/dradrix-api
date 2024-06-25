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
		"CheckOTP":  testCheckOTP,
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

type mockCheckOTPDB struct {
	DoesUserExistsFunc func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error)
	CheckOTPTxFunc     func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error)
}

func (mdb *mockCheckOTPDB) DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
	return mdb.DoesUserExistsFunc(ctx, arg)
}

func (mdb *mockCheckOTPDB) CheckOTPTx(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
	return mdb.CheckOTPTxFunc(ctx, arg)
}

func testCheckOTP(t *testing.T, appHandler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return nil, nil
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CheckOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("error user not exists", func(t *testing.T) {
		dataRequest := handlers.CheckOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
			PinCode:     "0123",
		}
		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return nil, errors.New("error does not exists")
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CHK_OTP_01"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CheckOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("users does not exists", func(t *testing.T) {
		dataRequest := handlers.CheckOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
			PinCode:     "0123",
		}

		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return nil, nil
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CHK_OTP_02"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CheckOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("error when checking otp", func(t *testing.T) {
		dataRequest := handlers.CheckOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
			PinCode:     "0123",
		}
		user := models.User{
			PhoneNumber: dataRequest.PhoneNumber,
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),
			Email:       gofaker.Email(),
		}

		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return &user, nil
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return nil, errors.New("error when checking otp")
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CHK_OTP_03"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CheckOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("otp does not exists", func(t *testing.T) {
		dataRequest := handlers.CheckOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
			PinCode:     "0123",
		}
		user := models.User{
			PhoneNumber: dataRequest.PhoneNumber,
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),
			Email:       gofaker.Email(),
		}

		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return &user, nil
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return nil, nil
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		wantStatus := http.StatusBadRequest
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}
		wantCode := "ERR_AUTH_CHK_OTP_04"
		if !strings.HasPrefix(response, wantCode) {
			t.Fatalf("CheckOTP(): response error - got %s, want %s", response, wantCode)
		}
	})

	t.Run("success", func(t *testing.T) {
		dataRequest := handlers.CheckOTPRequest{
			PhoneNumber: gofaker.Phonenumber(),
			Language:    "en",
			PinCode:     "0123",
		}
		user := models.User{
			PhoneNumber: dataRequest.PhoneNumber,
			FirstName:   gofaker.FirstName(),
			LastName:    gofaker.LastName(),
			Email:       gofaker.Email(),
		}
		otp := models.OTP{
			Id:          primitive.NewObjectID(),
			WaMessageId: "",
			PhoneNumber: dataRequest.PhoneNumber,
			PinCode:     dataRequest.PinCode,
			Active:      false,
		}

		mux := chi.NewMux()
		db := &mockCheckOTPDB{
			DoesUserExistsFunc: func(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
				return &user, nil
			},
			CheckOTPTxFunc: func(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
				return &otp, nil
			},
		}

		appHandler.CheckOTP(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/otp/check",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{},
		)
		// if user.PhoneNumber != dataRequest.PhoneNumber {
		// 	t.Fatalf("CheckOTP(): PhoneNumber - got %d; want %d", user.PhoneNumber, dataRequest.PhoneNumber)
		// }
		// if otp.PinCode != dataRequest.PinCode {
		// 	t.Fatalf("CheckOTP(): PinCode - got %d; want %d", otp.PinCode, dataRequest.PinCode)
		// }
		wantStatus := http.StatusOK
		if code != wantStatus {
			t.Fatalf("CheckOTP(): status - got %d; want %d", code, wantStatus)
		}

		got := handlers.CheckOTPResponse{}
		json.Unmarshal([]byte(response), &got)
		if got.User.PhoneNumber != user.PhoneNumber {
			t.Fatalf("CheckOTP(): PhoneNumber - got %s; want %s", got.User.PhoneNumber, dataRequest.PhoneNumber)
		}
	})
}
