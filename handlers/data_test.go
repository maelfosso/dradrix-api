package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
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
	Name:        gofaker.Name(),
	PhoneNumber: gofaker.Phonenumber(),
}

func TestData(t *testing.T) {
	handler := handlers.NewAppHandler()
	handler.GetAuthenticatedUser = func(r *http.Request) *models.User {
		return authenticatedUser
	}

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		// "DataMiddleware":   testDataMiddleware,
		// "GetAllActivities": testGetAllActivities,
		"CreateData": testCreateData,
		// "GetData":          testGetData,
		// "UpdateData":       testUpdateData,
		// "DeleteData":       testDeleteData,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockCreateDataDB struct {
	CreateDataFunc func(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error)
}

func (mdb *mockCreateDataDB) CreateData(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error) {
	return mdb.CreateDataFunc(ctx, arg)
}

func testCreateData(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateDataDB{
			CreateDataFunc: func(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error) {
				return nil, nil
			},
		}

		handler.CreateData(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateData(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("CreateData(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.App().Name(),
			Description: gofaker.Paragraph(),
		}
		dataRequest := handlers.CreateDataRequest{
			Values: map[string]any{
				"n_devis":    gofaker.UUIDHyphenated(),
				"n_os":       gofaker.UUIDDigit(),
				"date_os":    gofaker.Date(),
				"montant_os": sfaker.Number().Number(7),
			},
		}
		db := &mockCreateDataDB{
			CreateDataFunc: func(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.CreateData(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{
				{
					Name:  "activity",
					Value: activity,
				},
			},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateData(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_DATA_CRT_01"
		if response != want {
			t.Fatalf("CreateData(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.App().Name(),
			Description: gofaker.Paragraph(),
		}
		dataValues := map[string]any{
			"n_devis":    gofaker.UUIDHyphenated(),
			"n_os":       gofaker.UUIDDigit(),
			"date_os":    gofaker.Date(),
			"montant_os": sfaker.Number().Number(7),
		}
		data := &models.Data{
			Id:     primitive.NewObjectID(),
			Values: dataValues,

			ActivityId: activity.Id,
			CreatedBy:  authenticatedUser.Id,
		}

		mux := chi.NewMux()
		db := &mockCreateDataDB{
			CreateDataFunc: func(ctx context.Context, arg storage.CreateDataParams) (*models.Data, error) {
				return data, nil
			},
		}

		handler.CreateData(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateDataRequest{
				Values: dataValues,
			},
			[]helpertest.ContextData{{Name: "activity", Value: activity}},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("CreateData(): status - got %d; want %d", code, want)
		}

		got := handlers.CreateDataResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := dataEq(&got.Data, data); err != nil {
			t.Fatalf("CreateData(): %v", err)
		}
		if got.Data.ActivityId != activity.Id {
			t.Fatalf("CreateData(): activityId - got %s; want %s", got.Data.ActivityId, activity.Id)
		}
		if got.Data.CreatedBy != authenticatedUser.Id {
			t.Fatalf("CreateData(): CreatedBy - got %s; want %s", got.Data.CreatedBy, authenticatedUser.Id)
		}
	})

}

func dataEq(got, want *models.Data) error {
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
		return fmt.Errorf("Id - got = %s; want %s", got.Id, want.Id)
	}
	if len(got.Values) != len(want.Values) {
		return fmt.Errorf("#Values - got %d; want %d", len(got.Values), len(want.Values))
	}

	gotValues := got.Values
	wantValues := want.Values
	if !reflect.DeepEqual(gotValues, wantValues) {
		return fmt.Errorf("Values - got %+v; want %+v", gotValues, wantValues)
	}

	if got.CreatedBy != want.CreatedBy {
		return fmt.Errorf("CreatedBy - got %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	if got.ActivityId != want.ActivityId {
		return fmt.Errorf("ActivityId - got %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
