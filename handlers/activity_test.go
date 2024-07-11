package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func TestActivity(t *testing.T) {
	handler := handlers.NewAppHandler()

	tests := map[string]func(*testing.T, *handlers.AppHandler){
		// "ActivityMiddleware": testActivityMiddleware,
		// "GetAllActivities":   testGetAllActivities,
		// "CreateActivity":     testCreateActivity,
		// "GetActivity":        testGetActivity,
		"UpdateActivity": testUpdateActivity,
		// "DeleteActivity":     testDeleteActivity,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc(t, handler)
		})
	}
}

type mockAcitivityMiddlewareDB struct {
	GetActivityFunc func(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error)
}

func (mdb *mockAcitivityMiddlewareDB) GetActivity(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error) {
	return mdb.GetActivityFunc(ctx, arg)
}

func testActivityMiddleware(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid activity id", func(t *testing.T) {
		mux := chi.NewRouter()
		db := &mockAcitivityMiddlewareDB{}
		db.GetActivityFunc = func(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error) {
			return nil, nil
		}

		mux.Route("/{activityId}", func(r chi.Router) {
			handler.ActivityMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		_, w, response := helpertest.MakeGetRequest(mux, "/1", []helpertest.ContextData{})
		code := w.StatusCode
		wantStatusCode := http.StatusBadRequest
		if code != wantStatusCode {
			t.Fatalf("ActivityMiddleware(): status - got %d; want %d", code, wantStatusCode)
		}
		wantError := "ERR_ATVT_MDW_01"
		if response != wantError {
			t.Fatalf("ActivityMiddleware(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewRouter()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		db := &mockAcitivityMiddlewareDB{}
		db.GetActivityFunc = func(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error) {
			return nil, errors.New("error from db")
		}

		mux.Route("/{activityId}", func(r chi.Router) {
			handler.ActivityMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		_, w, response := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{
				{
					Name:  "organization",
					Value: organization,
				},
			})
		code := w.StatusCode
		wantStatusCode := http.StatusBadRequest
		if code != wantStatusCode {
			t.Fatalf("ActivityMiddleware(): status - got %d; want %d", code, wantStatusCode)
		}
		wantError := "ERR_ATVT_MDW_02"
		if response != wantError {
			t.Fatalf("ActivityMiddleware(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("no activity found", func(t *testing.T) {
		mux := chi.NewRouter()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		db := &mockAcitivityMiddlewareDB{}
		db.GetActivityFunc = func(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error) {
			return nil, nil
		}

		mux.Route("/{activityId}", func(r chi.Router) {
			handler.ActivityMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		_, w, response := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{
				{
					Name:  "organization",
					Value: organization,
				},
			})
		code := w.StatusCode
		wantStatusCode := http.StatusNotFound
		if code != wantStatusCode {
			t.Fatalf("ActivityMiddleware(): status - got %d; want %d", code, wantStatusCode)
		}
		wantError := "ERR_ATVT_MDW_03"
		if response != wantError {
			t.Fatalf("ActivityMiddleware(): response error - got %s; want %s", response, wantError)
		}
	})

	t.Run("activity found", func(t *testing.T) {
		mux := chi.NewRouter()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}
		db := &mockAcitivityMiddlewareDB{}
		db.GetActivityFunc = func(ctx context.Context, arg storage.GetActivityParams) (*models.Activity, error) {
			return activity, nil
		}

		mux.Route("/{activityId}", func(r chi.Router) {
			handler.ActivityMiddleware(r, db)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				gotActivity := r.Context().Value("activity").(*models.Activity)
				if err := activityEq(gotActivity, activity); err != nil {
					t.Fatalf("ActivityMiddleware(): %v", err)
				}
			})
		})

		_, w, _ := helpertest.MakeGetRequest(
			mux,
			fmt.Sprintf("/%s", primitive.NewObjectID().Hex()),
			[]helpertest.ContextData{
				{
					Name:  "organization",
					Value: organization,
				},
			})
		code := w.StatusCode
		wantStatusCode := http.StatusOK
		if code != wantStatusCode {
			t.Fatalf("ActivityMiddleware(): status - got %d; want %d", code, wantStatusCode)
		}
	})
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
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		mockDb.GetAllActivitiesFunc = func(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error) {
			return []*models.Activity{}, errors.New("error from db")
		}

		handler.GetAllActivities(mux, mockDb)
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
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("GetAllActivities(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_ATVT_GALL_01"
		if response != wantError {
			t.Fatalf("GetAllActivities(): status - got %s; want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}

		activities := []*models.Activity{
			{
				Name:        sfaker.Hacker().Noun(),
				Description: gofaker.Paragraph(),

				Fields: []models.ActivityField{
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "number",
						Key:         false,
					},
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "text",
					},
				},

				OrganizationId: organization.Id,
				CreatedBy:      primitive.NewObjectID(),
			},
		}

		mockDb.GetAllActivitiesFunc = func(ctx context.Context, arg storage.GetAllActivitiesParams) ([]*models.Activity, error) {
			return activities, nil
		}

		handler.GetAllActivities(mux, mockDb)
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
		wantCode := http.StatusOK
		if code != wantCode {
			t.Fatalf("GetAllActivities(): status - got %d; want %d", code, wantCode)
		}

		got := handlers.GetAllActivitiesResponse{}
		json.Unmarshal([]byte(response), &got)
		if len(got.Activities) != len(activities) {
			t.Fatalf("GetAllActivities(): Id - got %d; want %d", len(got.Activities), len(activities))
		}
		for i := 0; i < len(got.Activities); i++ {
			if err := activityEq(got.Activities[i], activities[i]); err != nil {
				t.Fatalf("GetAllActivities(): %v", err)
			}
		}
	})
}

type mockCreateActivityDB struct {
	CreateActivityFunc func(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error)
}

func (mdb *mockCreateActivityDB) CreateActivity(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error) {
	return mdb.CreateActivityFunc(ctx, arg)
}

func testCreateActivity(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockCreateActivityDB{
			CreateActivityFunc: func(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error) {
				return nil, nil
			},
		}

		handler.CreateActivity(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateActivity(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("CreateActivity(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		db := &mockCreateActivityDB{
			CreateActivityFunc: func(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.CreateActivity(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateActivityRequest{
				Name:        sfaker.Company().Name(),
				Description: gofaker.Paragraph(),
			},
			[]helpertest.ContextData{
				{
					Name:  "organization",
					Value: organization,
				},
			},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("CreateActivity(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_ATVT_CRT_01"
		if response != want {
			t.Fatalf("CreateActivity(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}

		mux := chi.NewMux()
		db := &mockCreateActivityDB{
			CreateActivityFunc: func(ctx context.Context, arg storage.CreateActivityParams) (*models.Activity, error) {
				return activity, nil
			},
		}

		handler.CreateActivity(mux, db)
		code, _, response := helpertest.MakePostRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.CreateActivityRequest{
				Name:        organization.Name,
				Description: organization.Bio,
			},
			[]helpertest.ContextData{{Name: "organization", Value: organization}},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("CreateActivity(): status - got %d; want %d", code, want)
		}

		got := handlers.CreateActivityResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := activityEq(&got.Activity, activity); err != nil {
			t.Fatalf("CreateActivity(): %v", err)
		}
		if got.Activity.OrganizationId != organization.Id {
			t.Fatalf("CreateActivity(): organizationId - got %s; want %s", got.Activity.OrganizationId, organization.Id)
		}
	})
}

func testGetActivity(t *testing.T, handler *handlers.AppHandler) {

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: primitive.NewObjectID(),
			CreatedBy:      primitive.NewObjectID(),
		}

		db := &struct{}{}

		handler.GetActivity(mux, db)
		_, w, response := helpertest.MakeGetRequest(
			mux,
			"/",
			[]helpertest.ContextData{
				{
					Name:  "activity",
					Value: activity,
				},
			},
		)
		code := w.StatusCode
		if code != http.StatusOK {
			t.Fatalf("GetActivity(): status - got %d; want %d", code, http.StatusOK)
		}

		got := handlers.GetActivityResponse{}
		json.Unmarshal([]byte(response), &got)
		if err := activityEq(&got.Activity, activity); err != nil {
			t.Fatalf("GetActivity(): %v", err)
		}
	})
}

type mockUpdateActivityDB struct {
	UpdateSetInActivityFunc      func(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error)
	UpdateAddToActivityFunc      func(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error)
	UpdateRemoveFromActivityFunc func(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error)
}

func (mdb *mockUpdateActivityDB) UpdateSetInActivity(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error) {
	return mdb.UpdateSetInActivityFunc(ctx, arg)
}

func (mdb *mockUpdateActivityDB) UpdateAddToActivity(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error) {
	return mdb.UpdateAddToActivityFunc(ctx, arg)
}

func (mdb *mockUpdateActivityDB) UpdateRemoveFromActivity(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error) {
	return mdb.UpdateRemoveFromActivityFunc(ctx, arg)
}

func testUpdateActivity(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateActivityDB{}

		handler.UpdateActivity(mux, db)
		code, _, response := helpertest.MakePatchRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			"{\"test\": \"that\"}",
			[]helpertest.ContextData{},
		)
		if code != http.StatusBadRequest {
			t.Fatalf("UpdateActivity(): status - got %d; want %d", code, http.StatusBadRequest)
		}
		want := "ERR_HDL_PRB_"
		if !strings.HasPrefix(response, want) {
			t.Fatalf("UpdateActivity(): response error - got %s, want %s", response, want)
		}
	})

	t.Run("wrong input.field value", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}
		db := &mockUpdateActivityDB{
			UpdateSetInActivityFunc: func(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error) {
				return nil, nil
			},
			UpdateAddToActivityFunc: func(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error) {
				return nil, nil
			},
			UpdateRemoveFromActivityFunc: func(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error) {
				return nil, nil
			},
		}

		testCases := map[string]struct {
			Input          handlers.UpdateActivityRequest
			HttpStatusCode int
			ResponseError  string
		}{
			"set": {
				Input: handlers.UpdateActivityRequest{
					Operation: "set",
					Field:     "",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_010",
			},
			"add": {
				Input: handlers.UpdateActivityRequest{
					Operation: "add",
					Field:     "Field",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_011",
			},
			"remove": {
				Input: handlers.UpdateActivityRequest{
					Operation: "remove",
					Field:     "X",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_012",
			},
			"something else": {
				Input: handlers.UpdateActivityRequest{
					Operation: "sth",
					Field:     "Name",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_013",
			},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				handler.UpdateActivity(mux, db)
				code, _, response := helpertest.MakePatchRequest(
					mux,
					"/",
					helpertest.CreateFormHeader(),
					tc.Input,
					[]helpertest.ContextData{
						{Name: "organization", Value: organization},
						{Name: "activity", Value: activity},
					},
				)
				wantCode := tc.HttpStatusCode
				if code != wantCode {
					t.Fatalf("UpdateActivity(): %s - status - got %d; want %d", name, code, wantCode)
				}
				wantError := tc.ResponseError
				if response != wantError {
					t.Fatalf("UpdateActivity(): %s - response error - got %s, want %s", name, response, wantError)
				}
			})
		}
	})

	// t.Run("wrong input.value value", func(t *testing.T) {
	// 	mux := chi.NewMux()
	// 	organization := &models.Organization{
	// 		Id:          primitive.NewObjectID(),
	// 		Name:        sfaker.Company().Name(),
	// 		Description: gofaker.Paragraph(),
	// 	}
	// 	activity := &models.Activity{
	// 		Id:          primitive.NewObjectID(),
	// 		Name:        sfaker.Hacker().Noun(),
	// 		Description: gofaker.Paragraph(),

	// 		Fields: []models.ActivityFields{
	// 			{
	// 				Code:        sfaker.App().Name(),
	// 				Name:        sfaker.App().String(),
	// 				Description: gofaker.Paragraph(),
	// 				Type:        "number",
	// 				Id:          false,
	// 			},
	// 			{
	// 				Code:        sfaker.App().Name(),
	// 				Name:        sfaker.App().String(),
	// 				Description: gofaker.Paragraph(),
	// 				Type:        "text",
	// 			},
	// 		},

	// 		OrganizationId: organization.Id,
	// 		CreatedBy: primitive.NewObjectID(),
	// 	}
	// 	db := &mockUpdateActivityDB{
	// 		UpdateSetInActivityFunc: func(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error) {
	// 			return nil, nil
	// 		},
	// 		UpdateAddToActivityFunc: func(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error) {
	// 			return nil, nil
	// 		},
	// 		UpdateRemoveFromActivityFunc: func(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error) {
	// 			return nil, nil
	// 		},
	// 	}

	// 	testCases := map[string]struct {
	// 		Input          handlers.UpdateActivityRequest
	// 		HttpStatusCode int
	// 		ResponseError  string
	// 	}{
	// 		// "set": {
	// 		// 	Input: handlers.UpdateActivityRequest{
	// 		// 		Operation: "set",
	// 		// 		Field:     "",
	// 		// 		Value:     sfaker.App().String(),
	// 		// 	},
	// 		// 	HttpStatusCode: http.StatusBadRequest,
	// 		// 	ResponseError:  "ERR_ATVT_UDT_010",
	// 		// },
	// 		"add": {
	// 			Input: handlers.UpdateActivityRequest{
	// 				Operation: "add",
	// 				Field:     "Fields",
	// 				Value:     sfaker.App().String(),
	// 			},
	// 			HttpStatusCode: http.StatusBadRequest,
	// 			ResponseError:  "ERR_ATVT_UDT_014",
	// 		},
	// 		// "remove": {
	// 		// 	Input: handlers.UpdateActivityRequest{
	// 		// 		Operation: "remove",
	// 		// 		Field:     "X",
	// 		// 		Value:     sfaker.App().String(),
	// 		// 	},
	// 		// 	HttpStatusCode: http.StatusBadRequest,
	// 		// 	ResponseError:  "ERR_ATVT_UDT_012",
	// 		// },
	// 		// "something else": {
	// 		// 	Input: handlers.UpdateActivityRequest{
	// 		// 		Operation: "sth",
	// 		// 		Field:     "Name",
	// 		// 		Value:     sfaker.App().String(),
	// 		// 	},
	// 		// 	HttpStatusCode: http.StatusBadRequest,
	// 		// 	ResponseError:  "ERR_ATVT_UDT_013",
	// 		// },
	// 	}

	// 	for name, tc := range testCases {
	// 		t.Run(name, func(t *testing.T) {
	// 			handler.UpdateActivity(mux, db)
	// 			code, _, response := helpertest.MakePatchRequest(
	// 				mux,
	// 				"/",
	// 				helpertest.CreateFormHeader(),
	// 				tc.Input,
	// 				[]helpertest.ContextData{
	// 					{Name: "organization", Value: organization},
	// 					{Name: "activity", Value: activity},
	// 				},
	// 			)
	// 			wantCode := tc.HttpStatusCode
	// 			if code != wantCode {
	// 				t.Fatalf("UpdateActivity(): %s - status - got %d; want %d", name, code, wantCode)
	// 			}
	// 			wantError := tc.ResponseError
	// 			if response != wantError {
	// 				t.Fatalf("UpdateActivity(): %s - response error - got %s, want %s", name, response, wantError)
	// 			}
	// 		})
	// 	}
	// })

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}
		db := &mockUpdateActivityDB{
			UpdateSetInActivityFunc: func(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error) {
				return nil, errors.New("an error happens")
			},
			UpdateAddToActivityFunc: func(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error) {
				return nil, errors.New("an error happens")
			},
			UpdateRemoveFromActivityFunc: func(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error) {
				return nil, errors.New("an error happens")
			},
		}

		testCases := map[string]struct {
			Input          handlers.UpdateActivityRequest
			HttpStatusCode int
			ResponseError  string
		}{
			"set": {
				Input: handlers.UpdateActivityRequest{
					Operation: "set",
					Field:     "Name",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_02",
			},
			"add": {
				Input: handlers.UpdateActivityRequest{
					Operation: "add",
					Field:     "Fields",
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_02",
			},
			"remove": {
				Input: handlers.UpdateActivityRequest{
					Operation: "remove",
					Field:     "Fields",
					Position:  1,
					Value:     sfaker.App().String(),
				},
				HttpStatusCode: http.StatusBadRequest,
				ResponseError:  "ERR_ATVT_UDT_02",
			},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				handler.UpdateActivity(mux, db)
				code, _, response := helpertest.MakePatchRequest(
					mux,
					"/",
					helpertest.CreateFormHeader(),
					tc.Input,
					[]helpertest.ContextData{
						{Name: "organization", Value: organization},
						{Name: "activity", Value: activity},
					},
				)
				wantCode := tc.HttpStatusCode
				if code != wantCode {
					t.Fatalf("UpdateActivity(): %s - status - got %d; want %d", name, code, wantCode)
				}
				wantError := tc.ResponseError
				if response != wantError {
					t.Fatalf("UpdateActivity(): %s - response error - got %s, want %s", name, response, wantError)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		var organization *models.Organization
		var activity *models.Activity

		dataRequest := []handlers.UpdateActivityRequest{
			{
				Operation: "set",
				Field:     "Name",
				Value:     sfaker.App().String(),
			},
			{
				Operation: "add",
				Field:     "Fields",
				Value: models.ActivityField{
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Code:        sfaker.App().Name(),
				},
				Position: 0,
			},
			{
				Operation: "remove",
				Field:     "Fields",
				Position:  1,
				// Value:     sfaker.App().String(),
			},
			{
				Operation: "sth",
				Field:     "Name",
				Value:     sfaker.App().String(),
			},
		}

		mux := chi.NewMux()
		db := &mockUpdateActivityDB{
			UpdateSetInActivityFunc: func(ctx context.Context, arg storage.UpdateSetInActivityParams) (*models.Activity, error) {
				updatedActivity := activity
				updatedActivity.Name = dataRequest[0].Value.(string)
				return updatedActivity, nil
			},
			UpdateAddToActivityFunc: func(ctx context.Context, arg storage.UpdateAddToActivityParams) (*models.Activity, error) {
				updatedActivity := &models.Activity{
					Id:             activity.Id,
					Name:           activity.Name,
					Description:    activity.Description,
					OrganizationId: activity.OrganizationId,
					CreatedBy:      activity.CreatedBy,
					Fields:         []models.ActivityField{},
				}
				updatedActivity.Fields = append(updatedActivity.Fields, activity.Fields...)
				updatedActivity.Fields = append(updatedActivity.Fields, dataRequest[1].Value.(models.ActivityField))
				return updatedActivity, nil
			},
			UpdateRemoveFromActivityFunc: func(ctx context.Context, arg storage.UpdateRemoveFromActivityParams) (*models.Activity, error) {
				updatedActivity := &models.Activity{
					Id:             activity.Id,
					Name:           activity.Name,
					Description:    activity.Description,
					OrganizationId: activity.OrganizationId,
					CreatedBy:      activity.CreatedBy,
					Fields:         []models.ActivityField{},
				}
				updatedActivity.Fields = []models.ActivityField{activity.Fields[0]}
				return updatedActivity, nil
			},
		}

		testCases := map[string]struct {
			Input          handlers.UpdateActivityRequest
			HttpStatusCode int
			CheckResponse  func(string, string)
		}{
			"set": {
				Input:          dataRequest[0],
				HttpStatusCode: http.StatusOK,
				CheckResponse: func(name, response string) {
					got := handlers.UpdateActivityResponse{}
					json.Unmarshal([]byte(response), &got)

					if got.Activity.Id != activity.Id {
						t.Fatalf("UpdateActivity(): Id - got %s; want %s", got.Activity.Id, organization.Id)
					}
					if got.Activity.OrganizationId != organization.Id {
						t.Fatalf("UpdateActivity(): OrganizationId - got %s; want %s", got.Activity.Id, organization.Id)
					}
					if got.Activity.Name != dataRequest[0].Value.(string) {
						t.Fatalf("UpdateActivity(): %s - Name - got %s; want %s", name, got.Activity.Name, dataRequest[0].Value)
					}
				},
			},
			"add": {
				Input:          dataRequest[1],
				HttpStatusCode: http.StatusOK,
				CheckResponse: func(name, response string) {
					got := handlers.UpdateActivityResponse{}
					json.Unmarshal([]byte(response), &got)

					if got.Activity.Id != activity.Id {
						t.Fatalf("UpdateActivity(): %s - Id - got %s; want %s", name, got.Activity.Id, organization.Id)
					}
					if got.Activity.OrganizationId != organization.Id {
						t.Fatalf("UpdateActivity(): %s - OrganizationId - got %s; want %s", name, got.Activity.Id, organization.Id)
					}
					if len(got.Activity.Fields) != len(activity.Fields)+1 {
						t.Fatalf("UpdateActivity(): %s - Length Fields - got %d; want %d", name, len(got.Activity.Fields), len(activity.Fields)+1)
					}
					if got.Activity.Fields[len(got.Activity.Fields)-1] != dataRequest[1].Value.(models.ActivityField) {
						t.Fatalf("UpdateActivity(): %s - Last Fields - got %+v; want %+v", name, got.Activity.Fields[len(got.Activity.Fields)-1], dataRequest[1].Value.(models.ActivityField))
					}
				},
			},
			"remove": {
				Input:          dataRequest[2],
				HttpStatusCode: http.StatusOK,
				CheckResponse: func(name, response string) {
					got := handlers.UpdateActivityResponse{}
					json.Unmarshal([]byte(response), &got)

					if got.Activity.Id != activity.Id {
						t.Fatalf("UpdateActivity(): %s - Id - got %s; want %s", name, got.Activity.Id, organization.Id)
					}
					if got.Activity.OrganizationId != organization.Id {
						t.Fatalf("UpdateActivity(): %s - OrganizationId - got %s; want %s", name, got.Activity.Id, organization.Id)
					}
					if len(got.Activity.Fields) != len(activity.Fields)-1 {
						t.Fatalf("UpdateActivity(): %s - Length Fields - got %d; want %d", name, len(got.Activity.Fields), len(activity.Fields)-1)
					}
					if got.Activity.Fields[0] != activity.Fields[0] {
						t.Fatalf("UpdateActivity(): %s - Last Fields - got %v; want %v", name, got.Activity.Fields[len(got.Activity.Fields)-1], dataRequest[1].Value.(models.ActivityField))
					}
				},
			},
		}

		for name, tc := range testCases {
			organization = &models.Organization{
				Id:   primitive.NewObjectID(),
				Name: sfaker.Company().Name(),
				Bio:  gofaker.Paragraph(),
			}
			activity = &models.Activity{
				Id:          primitive.NewObjectID(),
				Name:        sfaker.Hacker().Noun(),
				Description: gofaker.Paragraph(),

				Fields: []models.ActivityField{
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "number",
						Key:         true,
					},
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "text",
					},
				},

				OrganizationId: organization.Id,
				CreatedBy:      primitive.NewObjectID(),
			}

			t.Run(name, func(t *testing.T) {
				handler.UpdateActivity(mux, db)
				code, _, response := helpertest.MakePatchRequest(
					mux,
					"/",
					helpertest.CreateFormHeader(),
					tc.Input,
					[]helpertest.ContextData{
						{Name: "organization", Value: organization},
						{Name: "activity", Value: activity},
					},
				)
				wantCode := tc.HttpStatusCode
				if code != wantCode {
					log.Println("Response : ", response)
					t.Fatalf("UpdateActivity(): %s - status - got %d; want %d", name, code, wantCode)
				}
				tc.CheckResponse(name, response)
			})
		}

		// handler.UpdateActivity(mux, db)
		// code, _, response := helpertest.MakePatchRequest(
		// 	mux,
		// 	"/",
		// 	helpertest.CreateFormHeader(),
		// 	dataRequest,
		// 	[]helpertest.ContextData{
		// 		{Name: "organization", Value: organization},
		// 		{Name: "activity", Value: activity},
		// 	},
		// )
		// want := http.StatusOK
		// if code != want {
		// 	t.Fatalf("UpdateActivity(): status - got %d; want %d", code, want)
		// }

		// got := handlers.UpdateActivityResponse{}
		// json.Unmarshal([]byte(response), &got)
		// if got.Activity.OrganizationId != organization.Id {
		// 	t.Fatalf("UpdateActivity(): Id - got %s; want %s", got.Activity.Id, organization.Id)
		// }
		// if got.Activity.Name != dataRequest.Value.(string) {
		// 	t.Fatalf("UpdateActivity(): Name - got %s; want %s", got.Activity.Name, dataRequest.Value)
		// }
	})
}

type mockDeleteActivityDB struct {
	DeleteActivityFunc func(ctx context.Context, arg storage.DeleteActivityParams) error
}

func (mdb *mockDeleteActivityDB) DeleteActivity(ctx context.Context, arg storage.DeleteActivityParams) error {
	return mdb.DeleteActivityFunc(ctx, arg)
}

func testDeleteActivity(t *testing.T, handler *handlers.AppHandler) {
	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}
		db := &mockDeleteActivityDB{
			DeleteActivityFunc: func(ctx context.Context, arg storage.DeleteActivityParams) error {
				return errors.New("an error happens")
			},
		}

		handler.DeleteActivity(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			nil,
			[]helpertest.ContextData{
				{Name: "organization", Value: organization},
				{Name: "activity", Value: activity},
			},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("DeleteActivity(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_ATVT_DLT_01"
		if response != wantError {
			t.Fatalf("DeleteActivity(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		mux := chi.NewMux()
		organization := &models.Organization{
			Id:   primitive.NewObjectID(),
			Name: sfaker.Company().Name(),
			Bio:  gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Key:         false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			OrganizationId: organization.Id,
			CreatedBy:      primitive.NewObjectID(),
		}
		db := &mockDeleteActivityDB{
			DeleteActivityFunc: func(ctx context.Context, arg storage.DeleteActivityParams) error {
				return nil
			},
		}

		handler.DeleteActivity(mux, db)
		code, _, response := helpertest.MakeDeleteRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			nil,
			[]helpertest.ContextData{
				{Name: "organization", Value: organization},
				{Name: "activity", Value: activity},
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("DeleteActivity(): status - got %d; want %d", code, want)
		}

		got := handlers.DeleteActivityResponse{}
		json.Unmarshal([]byte(response), &got)
		if !got.Deleted {
			t.Fatalf("DeleteActivity(): deleted - got %v; want %v", got.Deleted, true)
		}
	})
}

func activityEq(got, want *models.Activity) error {
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
	if len(got.Fields) != len(want.Fields) {
		return fmt.Errorf("got.#Fields = %d; want %d", len(got.Fields), len(want.Fields))
	}

	gotFields := got.Fields
	sort.Slice(gotFields, func(i, j int) bool {
		return gotFields[i].Code < gotFields[j].Code
	})
	wantFields := want.Fields
	sort.Slice(wantFields, func(i, j int) bool {
		return wantFields[i].Code < wantFields[j].Code
	})
	n := len(gotFields)
	for i := 0; i < n; i++ {
		g := gotFields[i]
		w := wantFields[i]

		if g.Key != w.Key ||
			g.Code != w.Code ||
			g.Name != w.Name ||
			g.Description != w.Description ||
			g.Type != w.Type {

			return fmt.Errorf("got.Fields[%d] = %+v; want.Fields[%d] = %+v", i, g, i, w)
		}
	}

	if got.CreatedBy != want.CreatedBy {
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
