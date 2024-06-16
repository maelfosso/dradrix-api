package storage_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"

	sfaker "syreclabs.com/go/faker"
)

func TestData(t *testing.T) {
	integrationtest.SkipifShort(t)

	db, disconnect := integrationtest.CreateDatabase()
	defer disconnect()

	tests := map[string]func(*testing.T, *storage.Database){
		"CreateData": testCreateData,
		// "DeleteData":   testDeleteData,
		// "UpdateData":   testUpdateData,
		// "GetAllActivities": testGetAllActivities,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			integrationtest.Cleanup(db)

			tc(t, db)
		})
	}
}

func testCreateData(t *testing.T, db *storage.Database) {

	beforeCount, err := db.GetCollection("datas").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() Before; err = %v; want nil", err)
	}

	arg := storage.CreateDataParams{
		Values: map[string]any{
			"n_devis":    faker.UUIDHyphenated(),
			"n_os":       faker.UUIDDigit(),
			"date_os":    faker.Date(),
			"montant_os": sfaker.Number().Number(7),
		},

		ActivityId: primitive.NewObjectID(),
		CreatedBy:  primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateData(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateData() err = %v; want nil", err)
	}
	if activity.Id.IsZero() {
		t.Fatalf("CreateData(): Id is nil; want non nil")
	}

	afterCount, err := db.GetCollection("datas").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() After; err = %v; want nil", err)
	}
	if afterCount-beforeCount != 1 {
		t.Fatalf("AfterCount - BeforeCount = %d; want = %d", afterCount-beforeCount, 1)
	}

	got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
		Id:         activity.Id,
		ActivityId: arg.ActivityId,
	})
	if err != nil {
		t.Fatalf("GetData(): got error %+v; want nit", err.Error())
	}
	if err := dataEq(got, activity); err != nil {
		t.Fatalf("GetData(): %v", err.Error())
	}
}

// func testDeleteData(t *testing.T, db *storage.Database) {
// 	arg := storage.CreateDataParams{
// 		Name:        "a1",
// 		Description: "Data 1",
// 		Values: []models.DataValues{
// 			{Name: "f1", Description: "Description 1", Type: "number"},
// 			{Name: "f2", Description: "Description 2", Type: "text"},
// 		},

// 		ActivityId: primitive.NewObjectID(),
// 		CreatedBy: primitive.NewObjectID(),
// 	}
// 	activity, err := db.Storage.CreateData(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("CreateData(): err = %v; want nil", err)
// 	}

// 	_, err = db.Storage.GetData(context.Background(), storage.GetDataParams{
// 		Id:        activity.Id,
// 		ActivityId: arg.ActivityId,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetData(): got error = %v; want nil", err)
// 	}

// 	err = db.Storage.DeleteData(context.Background(), storage.DeleteDataParams{
// 		Id:        activity.Id,
// 		ActivityId: arg.ActivityId,
// 	})
// 	if err != nil {
// 		t.Fatalf("DeleteDataFromCompany(): got err = %v; want nil", err)
// 	}

// 	activity, err = db.Storage.GetData(context.Background(), storage.GetDataParams{
// 		Id:        activity.Id,
// 		ActivityId: arg.ActivityId,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetData(): got error = %v; want nil", err)
// 	}
// 	if activity != nil {
// 		t.Fatalf("GetData(): got %v; want nil", activity)
// 	}

// 	datas, err := db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
// 		ActivityId: arg.ActivityId,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivitiesFromCompany(): got err = %v; want nil", err)
// 	}
// 	if len(datas) != 0 {
// 		t.Fatalf("GetAllActivitiesFromCompany(): got %d number of datas; want = 0", len(datas))
// 	}
// }

// func testUpdateData(t *testing.T, db *storage.Database) {
// 	arg := storage.CreateDataParams{
// 		Name:        "a1",
// 		Description: "Data 1",
// 		Values: []models.DataValues{
// 			{Name: "f1", Description: "Description 1", Type: "number"},
// 			{Name: "f2", Description: "Description 2", Type: "text", Id: true},
// 		},

// 		ActivityId: primitive.NewObjectID(),
// 		CreatedBy: primitive.NewObjectID(),
// 	}
// 	activity, err := db.Storage.CreateData(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("CreateData() err = %v; want nil", err)
// 	}
// 	if activity.Id.IsZero() {
// 		t.Fatalf("CreateData(): Id is nil; want non nil")
// 	}

// 	t.Run("set", func(t *testing.T) {
// 		argForUpdate := storage.UpdateSetInDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,

// 			Field: "name",
// 			Value: "a2",
// 		}
// 		updated, err := db.Storage.UpdateSetInData(context.Background(), argForUpdate)
// 		if err != nil {
// 			t.Fatalf("UpdateData(): got error %v; want nil", err)
// 		}
// 		if updated.Name != argForUpdate.Value {
// 			t.Fatalf(
// 				"UpdateData(): updated %s value - got %s; want %s",
// 				argForUpdate.Field, updated.Name, argForUpdate.Value,
// 			)
// 		}

// 		argForUpdate = storage.UpdateSetInDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,

// 			Field: "fields.1.code",
// 			Value: "f1",
// 		}
// 		updated, err = db.Storage.UpdateSetInData(context.Background(), argForUpdate)
// 		if err != nil {
// 			t.Fatalf("UpdateData(): got error %v; want nil", err)
// 		}
// 		if updated.Values[1].Code != argForUpdate.Value {
// 			t.Fatalf(
// 				"UpdateData(): updated %s value - got %s; want %s",
// 				argForUpdate.Field, updated.Name, argForUpdate.Value,
// 			)
// 		}

// 		argForUpdate = storage.UpdateSetInDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,

// 			Field: "fields.1.id",
// 			Value: false,
// 		}
// 		updated, err = db.Storage.UpdateSetInData(context.Background(), argForUpdate)
// 		if err != nil {
// 			t.Fatalf("UpdateData(): got error %v; want nil", err)
// 		}
// 		if updated.Values[0].Id != argForUpdate.Value {
// 			t.Fatalf(
// 				"UpdateData(): updated %s value - got %s; want %s",
// 				argForUpdate.Field, updated.Name, argForUpdate.Value,
// 			)
// 		}

// 		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,
// 		})
// 		if err != nil {
// 			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
// 		}
// 		if err := dataEq(got, updated); err != nil {
// 			t.Fatalf("GetData(): %v", err.Error())
// 		}
// 	})

// 	t.Run("add", func(t *testing.T) {
// 		argForUpdate := storage.UpdateAddToDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,

// 			Field: "fields",
// 			Value: models.DataValues{
// 				Name:        sfaker.App().String(),
// 				Description: gofaker.Paragraph(),
// 				Type:        "number",
// 				Code:        sfaker.App().Name(),
// 			},
// 			Position: rand.UintN(uint(len(activity.Values))),
// 		}
// 		updated, err := db.Storage.UpdateAddToData(context.Background(), argForUpdate)
// 		if err != nil {
// 			t.Fatalf("UpdateAddToData(): got error %v; want nil", err)
// 		}
// 		if updated.Values[argForUpdate.Position] != argForUpdate.Value {
// 			t.Fatalf(
// 				"UpdateAddToData(): updated %s value - got %+v; want %+v",
// 				argForUpdate.Field, updated.Values[argForUpdate.Position], argForUpdate.Value,
// 			)
// 		}

// 		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,
// 		})
// 		if err != nil {
// 			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
// 		}
// 		if err = dataEq(got, updated); err != nil {
// 			t.Fatalf("GetData(): %v", err.Error())
// 		}

// 		activity = updated
// 	})

// 	t.Run("remove", func(t *testing.T) {
// 		activity, _ = db.Storage.GetData(context.Background(), storage.GetDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,
// 		})
// 		argForUpdate := storage.UpdateRemoveFromDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,

// 			Field:    "fields",
// 			Position: rand.UintN(uint(len(activity.Values))),
// 		}
// 		updated, err := db.Storage.UpdateRemoveFromData(context.Background(), argForUpdate)
// 		if err != nil {
// 			t.Fatalf("UpdateRemoveFromData(): got error %v; want nil", err)
// 		}
// 		if len(updated.Values) != len(activity.Values)-1 {
// 			t.Fatalf(
// 				"UpdateRemoveFromData(): len fields - got %d; want %d",
// 				len(updated.Values), len(activity.Values)-1)
// 		}

// 		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
// 			Id:        activity.Id,
// 			ActivityId: arg.ActivityId,
// 		})
// 		if err != nil {
// 			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
// 		}
// 		if err = dataEq(got, updated); err != nil {
// 			t.Fatalf("GetData(): %v", err.Error())
// 		}
// 	})

// }

// func testGetAllActivities(t *testing.T, db *storage.Database) {
// 	const NUM_ACTIVITIES_CREATED = 3
// 	companyA := primitive.NewObjectID()
// 	companyB := primitive.NewObjectID()

// 	var datas []*models.Data

// 	for i := 0; i < NUM_ACTIVITIES_CREATED; i++ {
// 		activity, _ := db.Storage.CreateData(context.Background(), storage.CreateDataParams{
// 			Name:        sfaker.Company().Name(),
// 			Description: gofaker.Paragraph(),

// 			Values: []models.DataValues{
// 				{Code: sfaker.App().Name(), Name: gofaker.Name(), Description: gofaker.Paragraph(), Type: "number"},
// 				{Code: sfaker.App().Name(), Name: gofaker.Name(), Description: gofaker.Paragraph(), Type: "text", Id: true},
// 			},

// 			ActivityId: companyA,
// 			CreatedBy: primitive.NewObjectID(),
// 		})
// 		datas = append(datas, activity)
// 	}

// 	got, err := db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
// 		ActivityId: companyA,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivities(): got error; want nil")
// 	}
// 	if len(got) != NUM_ACTIVITIES_CREATED {
// 		t.Fatalf("GetAllActivities(): got %d datas; want %d datas", len(got), NUM_ACTIVITIES_CREATED)
// 	}
// 	if err := dataEq(got[0], datas[0]); err != nil {
// 		t.Fatalf("GetAllActivities(): %v", err)
// 	}
// 	if err := dataEq(got[1], datas[1]); err != nil {
// 		t.Fatalf("GetAllActivities(): %v", err)
// 	}
// 	if err := dataEq(got[len(got)-1], datas[len(datas)-1]); err != nil {
// 		t.Fatalf("GetAllActivities(): %v", err)
// 	}

// 	got, err = db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
// 		ActivityId: companyB,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivities(): got error; want nil")
// 	}
// 	if len(got) != 0 {
// 		t.Fatalf("GetAllActivities(): got %d datas; want %d datas", len(got), 0)
// 	}

// 	db.Storage.DeleteData(context.Background(), storage.DeleteDataParams{
// 		Id:        datas[0].Id,
// 		ActivityId: companyA,
// 	})
// 	got, err = db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
// 		ActivityId: companyA,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivities(): got error; want nil")
// 	}
// 	if len(got) != NUM_ACTIVITIES_CREATED-1 {
// 		t.Fatalf("GetAllActivities(): got %d datas; want %d datas", len(got), NUM_ACTIVITIES_CREATED)
// 	}
// 	if err := dataEq(got[0], datas[1]); err != nil {
// 		t.Fatalf("GetAllActivities(): %v", err)
// 	}
// }

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
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
