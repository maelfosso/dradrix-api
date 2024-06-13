package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		"ActivityMiddleware": testActivityMiddleware,
		"GetAllActivities":   testGetAllActivities,
		"CreateActivity":     testCreateActivity,
		"GetActivity":        testGetActivity,
		"UpdateActivity":     testUpdateActivity,
		"DeleteActivity":     testDeleteActivity,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}

		activities := []*models.Activity{
			{
				Name:        sfaker.Hacker().Noun(),
				Description: gofaker.Paragraph(),

				Fields: []models.ActivityFields{
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "number",
						Id:          false,
					},
					{
						Code:        sfaker.App().Name(),
						Name:        sfaker.App().String(),
						Description: gofaker.Paragraph(),
						Type:        "text",
					},
				},

				CompanyId: company.Id,
				CreatedBy: primitive.NewObjectID(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
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
					Name:  "company",
					Value: company,
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
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
				Name:        company.Name,
				Description: company.Description,
			},
			[]helpertest.ContextData{{Name: "company", Value: company}},
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
		if got.Activity.CompanyId != company.Id {
			t.Fatalf("CreateActivity(): companyId - got %s; want %s", got.Activity.CompanyId, company.Id)
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

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: primitive.NewObjectID(),
			CreatedBy: primitive.NewObjectID(),
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
	UpdateActivityFunc func(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error)
}

func (mdb *mockUpdateActivityDB) UpdateActivity(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error) {
	return mdb.UpdateActivityFunc(ctx, arg)
}

func testUpdateActivity(t *testing.T, handler *handlers.AppHandler) {
	t.Run("invalid input data", func(t *testing.T) {
		mux := chi.NewMux()
		db := &mockUpdateActivityDB{
			UpdateActivityFunc: func(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error) {
				return nil, nil
			},
		}

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

	t.Run("error from db", func(t *testing.T) {
		mux := chi.NewMux()
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
		}
		db := &mockUpdateActivityDB{
			UpdateActivityFunc: func(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error) {
				return nil, errors.New("an error happens")
			},
		}

		handler.UpdateActivity(mux, db)
		code, _, response := helpertest.MakePatchRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			handlers.UpdateActivityRequest{
				Field: "Name",
				Value: sfaker.App().String(),
			},
			[]helpertest.ContextData{
				{Name: "company", Value: company},
				{Name: "activity", Value: activity},
			},
		)
		wantCode := http.StatusBadRequest
		if code != wantCode {
			t.Fatalf("UpdateActivity(): status - got %d; want %d", code, wantCode)
		}
		wantError := "ERR_ATVT_UDT_01"
		if response != wantError {
			t.Fatalf("UpdateActivity(): response error - got %s, want %s", response, wantError)
		}
	})

	t.Run("success", func(t *testing.T) {
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
		}
		dataRequest := handlers.UpdateActivityRequest{
			Field: "Name",
			Value: sfaker.App().String(),
		}

		mux := chi.NewMux()
		db := &mockUpdateActivityDB{
			UpdateActivityFunc: func(ctx context.Context, arg storage.UpdateActivityParams) (*models.Activity, error) {
				activity.Name = dataRequest.Value.(string)
				return activity, nil
			},
		}

		handler.UpdateActivity(mux, db)
		code, _, response := helpertest.MakePatchRequest(
			mux,
			"/",
			helpertest.CreateFormHeader(),
			dataRequest,
			[]helpertest.ContextData{
				{Name: "company", Value: company},
				{Name: "activity", Value: activity},
			},
		)
		want := http.StatusOK
		if code != want {
			t.Fatalf("UpdateActivity(): status - got %d; want %d", code, want)
		}

		got := handlers.UpdateActivityResponse{}
		json.Unmarshal([]byte(response), &got)
		if got.Activity.CompanyId != company.Id {
			t.Fatalf("UpdateActivity(): Id - got %s; want %s", got.Activity.Id, company.Id)
		}
		if got.Activity.Name != dataRequest.Value.(string) {
			t.Fatalf("UpdateActivity(): Name - got %s; want %s", got.Activity.Name, dataRequest.Value)
		}
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
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
				{Name: "company", Value: company},
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
		company := &models.Company{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),
		}
		activity := &models.Activity{
			Id:          primitive.NewObjectID(),
			Name:        sfaker.Hacker().Noun(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityFields{
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "number",
					Id:          false,
				},
				{
					Code:        sfaker.App().Name(),
					Name:        sfaker.App().String(),
					Description: gofaker.Paragraph(),
					Type:        "text",
				},
			},

			CompanyId: company.Id,
			CreatedBy: primitive.NewObjectID(),
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
				{Name: "company", Value: company},
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

		if g.Id != w.Id ||
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
