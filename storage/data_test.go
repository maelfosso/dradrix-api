package storage_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"reflect"
	"strconv"
	"strings"
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
		"DeleteData": testDeleteData,
		"UpdateData": testUpdateData,
		"GetAllData": testGetAllData,
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
	data, err := db.Storage.CreateData(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateData() err = %v; want nil", err)
	}
	if data.Id.IsZero() {
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
		Id:         data.Id,
		ActivityId: arg.ActivityId,
	})
	if err != nil {
		t.Fatalf("GetData(): got error %+v; want nit", err.Error())
	}
	if err := dataEq(got, data); err != nil {
		t.Fatalf("GetData(): %v", err.Error())
	}
}

func testDeleteData(t *testing.T, db *storage.Database) {
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
	data, err := db.Storage.CreateData(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateData(): err = %v; want nil", err)
	}

	err = db.Storage.DeleteData(context.Background(), storage.DeleteDataParams{
		Id:         data.Id,
		ActivityId: arg.ActivityId,
	})
	if err != nil {
		t.Fatalf("DeleteDataFromOrganization(): got err = %v; want nil", err)
	}

	data, err = db.Storage.GetData(context.Background(), storage.GetDataParams{
		Id:         data.Id,
		ActivityId: arg.ActivityId,
	})
	if err != nil {
		t.Fatalf("GetData(): got error = %v; want nil", err)
	}
	if data != nil {
		t.Fatalf("GetData(): got %v; want nil", data)
	}

	datas, err := db.Storage.GetAllData(context.Background(), storage.GetAllDataParams{
		ActivityId: arg.ActivityId,
	})
	if err != nil {
		t.Fatalf("GetAllData(): got err = %v; want nil", err)
	}
	if len(datas) != 0 {
		t.Fatalf("GetAllData(): got %d number of datas; want = 0", len(datas))
	}
	for _, v := range datas {
		if v.Id == data.Id {
			if v.DeletedAt == nil {
				t.Fatalf("GetAllData(): search created data: got deletedAt nil; want %v", data.DeletedAt)
			}
			return
		}
	}
}

func testUpdateData(t *testing.T, db *storage.Database) {
	arg := storage.CreateDataParams{
		Values: map[string]any{
			"n_devis":    faker.UUIDHyphenated(),
			"n_os":       faker.UUIDDigit(),
			"date_os":    faker.Date(),
			"montant_os": sfaker.Number().Number(7),
			"images": []string{
				sfaker.Company().Logo(),
				sfaker.Company().Logo(),
				sfaker.Company().Logo(),
				sfaker.Company().Logo(),
			},
		},

		ActivityId: primitive.NewObjectID(),
		CreatedBy:  primitive.NewObjectID(),
	}
	data, err := db.Storage.CreateData(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateData() err = %v; want nil", err)
	}
	if data.Id.IsZero() {
		t.Fatalf("CreateData(): Id is nil; want non nil")
	}

	s := reflect.ValueOf(arg.Values["images"])
	images := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		images[i] = v.(string)
	}

	t.Run("set", func(t *testing.T) {
		argForUpdate := storage.UpdateSetInDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,

			Field: "n_devis",
			Value: faker.UUIDHyphenated(),
		}
		updated, err := db.Storage.UpdateSetInData(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateData(): got error %v; want nil", err)
		}
		if updated.Values[argForUpdate.Field] != argForUpdate.Value {
			t.Fatalf(
				"UpdateData(): updated %s value - got %s; want %s",
				argForUpdate.Field, updated.Values[argForUpdate.Field], argForUpdate.Value,
			)
		}

		argForUpdate = storage.UpdateSetInDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,

			Field: "images.1",
			Value: sfaker.Avatar().String(),
		}
		updated, err = db.Storage.UpdateSetInData(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateData(): got error %v; want nil", err)
		}

		ss := strings.Split(argForUpdate.Field, ".")
		field := ss[0]
		position, _ := strconv.Atoi(ss[1])
		s := reflect.ValueOf(updated.Values[field])
		images := make([]string, s.Len())
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i).Interface()
			images[i] = v.(string)
		}
		if images[position] != argForUpdate.Value {
			t.Fatalf(
				"UpdateData(): updated %s value - got %s; want %s",
				argForUpdate.Field, (updated.Values[field].([]interface{}))[position], argForUpdate.Value,
			)
		}

		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,
		})
		if err != nil {
			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
		}
		if err := dataEq(got, updated); err != nil {
			t.Fatalf("GetData(): %v", err.Error())
		}
	})

	t.Run("add", func(t *testing.T) {
		argForUpdate := storage.UpdateAddToDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,

			Field:    "images",
			Value:    sfaker.Company().Logo(),
			Position: rand.UintN(uint(len(images))),
		}
		updated, err := db.Storage.UpdateAddToData(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateAddToData(): got error %v; want nil", err)
		}

		s = reflect.ValueOf(updated.Values["images"])
		images = make([]string, s.Len())
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i).Interface()
			images[i] = v.(string)
		}
		if images[argForUpdate.Position] != argForUpdate.Value {
			t.Fatalf(
				"UpdateAddToData(): updated %s value - got %+v; want %+v",
				argForUpdate.Field, images[argForUpdate.Position], argForUpdate.Value,
			)
		}

		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,
		})
		if err != nil {
			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
		}
		if err = dataEq(got, updated); err != nil {
			t.Fatalf("GetData(): %v", err.Error())
		}

		data = updated
	})

	t.Run("remove", func(t *testing.T) {
		data, _ = db.Storage.GetData(context.Background(), storage.GetDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,
		})
		argForUpdate := storage.UpdateRemoveFromDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,

			Field:    "images",
			Position: rand.UintN(uint(len(images))),
		}
		updated, err := db.Storage.UpdateRemoveFromData(context.Background(), argForUpdate)

		s = reflect.ValueOf(updated.Values["images"])
		imagesUpdated := make([]string, s.Len())
		for i := 0; i < s.Len(); i++ {
			v := s.Index(i).Interface()
			imagesUpdated[i] = v.(string)
		}
		if err != nil {
			t.Fatalf("UpdateRemoveFromData(): got error %v; want nil", err)
		}
		if len(imagesUpdated) != len(images)-1 {
			t.Fatalf(
				"UpdateRemoveFromData(): len fields - got %d; want %d",
				len(imagesUpdated), len(images)-1)
		}

		got, err := db.Storage.GetData(context.Background(), storage.GetDataParams{
			Id:         data.Id,
			ActivityId: arg.ActivityId,
		})
		if err != nil {
			t.Fatalf("GetData(): got error %+v; want nit", err.Error())
		}
		if err = dataEq(got, updated); err != nil {
			t.Fatalf("GetData(): %v", err.Error())
		}
	})

}

func testGetAllData(t *testing.T, db *storage.Database) {
	const NUM_DATA_CREATED = 3
	activityA := primitive.NewObjectID()
	activityB := primitive.NewObjectID()

	var datas []*models.Data

	for i := 0; i < NUM_DATA_CREATED; i++ {
		data, _ := db.Storage.CreateData(context.Background(), storage.CreateDataParams{
			Values: map[string]any{
				"n_devis":    faker.UUIDHyphenated(),
				"n_os":       faker.UUIDDigit(),
				"date_os":    faker.Date(),
				"montant_os": sfaker.Number().Number(7),
			},

			ActivityId: activityA,
			CreatedBy:  primitive.NewObjectID(),
		})
		datas = append(datas, data)
	}

	got, err := db.Storage.GetAllData(context.Background(), storage.GetAllDataParams{
		ActivityId: activityA,
	})
	if err != nil {
		t.Fatalf("GetAllData(): got error; want nil")
	}
	if len(got) != NUM_DATA_CREATED {
		t.Fatalf("GetAllData(): got %d datas; want %d datas", len(got), NUM_DATA_CREATED)
	}
	if err := dataEq(got[0], datas[0]); err != nil {
		t.Fatalf("GetAllData(): %v", err)
	}
	if err := dataEq(got[1], datas[1]); err != nil {
		t.Fatalf("GetAllData(): %v", err)
	}
	if err := dataEq(got[len(got)-1], datas[len(datas)-1]); err != nil {
		t.Fatalf("GetAllData(): %v", err)
	}

	got, err = db.Storage.GetAllData(context.Background(), storage.GetAllDataParams{
		ActivityId: activityB,
	})
	if err != nil {
		t.Fatalf("GetAllData(): got error; want nil")
	}
	if len(got) != 0 {
		t.Fatalf("GetAllData(): got %d datas; want %d datas", len(got), 0)
	}

	db.Storage.DeleteData(context.Background(), storage.DeleteDataParams{
		Id:         datas[0].Id,
		ActivityId: activityA,
	})
	got, err = db.Storage.GetAllData(context.Background(), storage.GetAllDataParams{
		ActivityId: activityA,
	})
	if err != nil {
		t.Fatalf("GetAllData(): got error; want nil")
	}
	if len(got) != NUM_DATA_CREATED-1 {
		t.Fatalf("GetAllData(): got %d datas; want %d datas", len(got), NUM_DATA_CREATED)
	}
	if err := dataEq(got[0], datas[1]); err != nil {
		t.Fatalf("GetAllData(): %v", err)
	}
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
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
