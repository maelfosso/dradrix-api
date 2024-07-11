package storage_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sort"
	"testing"

	gofaker "github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
	sfaker "syreclabs.com/go/faker"
)

func TestActivity(t *testing.T) {
	integrationtest.SkipifShort(t)

	db, disconnect := integrationtest.CreateDatabase()
	defer disconnect()

	tests := map[string]func(*testing.T, *storage.Database){
		"CreateActivity":   testCreateActivity,
		"DeleteActivity":   testDeleteActivity,
		"UpdateActivity":   testUpdateActivity,
		"GetAllActivities": testGetAllActivities,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			integrationtest.Cleanup(db)

			tc(t, db)
		})
	}
}

func testCreateActivity(t *testing.T, db *storage.Database) {

	beforeCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() Before; err = %v; want nil", err)
	}

	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityField{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text"},
		},

		OrganizationId: primitive.NewObjectID(),
		CreatedBy:      primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity() err = %v; want nil", err)
	}
	if activity.Id.IsZero() {
		t.Fatalf("CreateActivity(): Id is nil; want non nil")
	}

	afterCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() After; err = %v; want nil", err)
	}
	if afterCount-beforeCount != 1 {
		t.Fatalf("AfterCount - BeforeCount = %d; want = %d", afterCount-beforeCount, 1)
	}

	got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id:             activity.Id,
		OrganizationId: arg.OrganizationId,
	})
	if err != nil {
		t.Fatalf("GetActivity(): got error %+v; want nit", err.Error())
	}
	if err := activityEq(got, activity); err != nil {
		t.Fatalf("GetActivity(): %v", err.Error())
	}
}

func testDeleteActivity(t *testing.T, db *storage.Database) {
	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityField{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text"},
		},

		OrganizationId: primitive.NewObjectID(),
		CreatedBy:      primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity(): err = %v; want nil", err)
	}

	_, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id:             activity.Id,
		OrganizationId: arg.OrganizationId,
	})
	if err != nil {
		t.Fatalf("GetActivity(): got error = %v; want nil", err)
	}

	err = db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
		Id:             activity.Id,
		OrganizationId: arg.OrganizationId,
	})
	if err != nil {
		t.Fatalf("DeleteActivityFromOrganization(): got err = %v; want nil", err)
	}

	activity, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id:             activity.Id,
		OrganizationId: arg.OrganizationId,
	})
	if err != nil {
		t.Fatalf("GetActivity(): got error = %v; want nil", err)
	}
	if activity != nil {
		t.Fatalf("GetActivity(): got %v; want nil", activity)
	}

	activities, err := db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
		OrganizationId: arg.OrganizationId,
	})
	if err != nil {
		t.Fatalf("GetAllActivitiesFromOrganization(): got err = %v; want nil", err)
	}
	if len(activities) != 0 {
		t.Fatalf("GetAllActivitiesFromOrganization(): got %d number of activities; want = 0", len(activities))
	}
}

func testUpdateActivity(t *testing.T, db *storage.Database) {
	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityField{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text", Key: true},
		},

		OrganizationId: primitive.NewObjectID(),
		CreatedBy:      primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity() err = %v; want nil", err)
	}
	if activity.Id.IsZero() {
		t.Fatalf("CreateActivity(): Id is nil; want non nil")
	}

	t.Run("set", func(t *testing.T) {
		argForUpdate := storage.UpdateSetInActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,

			Field: "name",
			Value: "a2",
		}
		updated, err := db.Storage.UpdateSetInActivity(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateActivity(): got error %v; want nil", err)
		}
		if updated.Name != argForUpdate.Value {
			t.Fatalf(
				"UpdateActivity(): updated %s value - got %s; want %s",
				argForUpdate.Field, updated.Name, argForUpdate.Value,
			)
		}

		argForUpdate = storage.UpdateSetInActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,

			Field: "fields.1.code",
			Value: "f1",
		}
		updated, err = db.Storage.UpdateSetInActivity(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateActivity(): got error %v; want nil", err)
		}
		if updated.Fields[1].Code != argForUpdate.Value {
			t.Fatalf(
				"UpdateActivity(): updated %s value - got %s; want %s",
				argForUpdate.Field, updated.Name, argForUpdate.Value,
			)
		}

		argForUpdate = storage.UpdateSetInActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,

			Field: "fields.1.id",
			Value: false,
		}
		updated, err = db.Storage.UpdateSetInActivity(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateActivity(): got error %v; want nil", err)
		}
		if updated.Fields[0].Key != argForUpdate.Value {
			t.Fatalf(
				"UpdateActivity(): updated %s value - got %s; want %s",
				argForUpdate.Field, updated.Name, argForUpdate.Value,
			)
		}

		got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,
		})
		if err != nil {
			t.Fatalf("GetActivity(): got error %+v; want nit", err.Error())
		}
		if err := activityEq(got, updated); err != nil {
			t.Fatalf("GetActivity(): %v", err.Error())
		}
	})

	t.Run("add", func(t *testing.T) {
		argForUpdate := storage.UpdateAddToActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,

			Field: "fields",
			Value: models.ActivityField{
				Name:        sfaker.App().String(),
				Description: gofaker.Paragraph(),
				Type:        "number",
				Code:        sfaker.App().Name(),
			},
			Position: rand.UintN(uint(len(activity.Fields))),
		}
		updated, err := db.Storage.UpdateAddToActivity(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateAddToActivity(): got error %v; want nil", err)
		}
		if updated.Fields[argForUpdate.Position] != argForUpdate.Value {
			t.Fatalf(
				"UpdateAddToActivity(): updated %s value - got %+v; want %+v",
				argForUpdate.Field, updated.Fields[argForUpdate.Position], argForUpdate.Value,
			)
		}

		got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,
		})
		if err != nil {
			t.Fatalf("GetActivity(): got error %+v; want nit", err.Error())
		}
		if err = activityEq(got, updated); err != nil {
			t.Fatalf("GetActivity(): %v", err.Error())
		}

		activity = updated
	})

	t.Run("remove", func(t *testing.T) {
		activity, _ = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,
		})
		argForUpdate := storage.UpdateRemoveFromActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,

			Field:    "fields",
			Position: rand.UintN(uint(len(activity.Fields))),
		}
		updated, err := db.Storage.UpdateRemoveFromActivity(context.Background(), argForUpdate)
		if err != nil {
			t.Fatalf("UpdateRemoveFromActivity(): got error %v; want nil", err)
		}
		if len(updated.Fields) != len(activity.Fields)-1 {
			t.Fatalf(
				"UpdateRemoveFromActivity(): len fields - got %d; want %d",
				len(updated.Fields), len(activity.Fields)-1)
		}

		got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
			Id:             activity.Id,
			OrganizationId: arg.OrganizationId,
		})
		if err != nil {
			t.Fatalf("GetActivity(): got error %+v; want nit", err.Error())
		}
		if err = activityEq(got, updated); err != nil {
			t.Fatalf("GetActivity(): %v", err.Error())
		}
	})

}

func testGetAllActivities(t *testing.T, db *storage.Database) {
	const NUM_ACTIVITIES_CREATED = 3
	organizationA := primitive.NewObjectID()
	organizationB := primitive.NewObjectID()

	var activities []*models.Activity

	for i := 0; i < NUM_ACTIVITIES_CREATED; i++ {
		activity, _ := db.Storage.CreateActivity(context.Background(), storage.CreateActivityParams{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),

			Fields: []models.ActivityField{
				{Code: sfaker.App().Name(), Name: gofaker.Name(), Description: gofaker.Paragraph(), Type: "number"},
				{Code: sfaker.App().Name(), Name: gofaker.Name(), Description: gofaker.Paragraph(), Type: "text", Key: true},
			},

			OrganizationId: organizationA,
			CreatedBy:      primitive.NewObjectID(),
		})
		activities = append(activities, activity)
	}

	got, err := db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
		OrganizationId: organizationA,
	})
	if err != nil {
		t.Fatalf("GetAllActivities(): got error; want nil")
	}
	if len(got) != NUM_ACTIVITIES_CREATED {
		t.Fatalf("GetAllActivities(): got %d activities; want %d activities", len(got), NUM_ACTIVITIES_CREATED)
	}
	if err := activityEq(got[0], activities[0]); err != nil {
		t.Fatalf("GetAllActivities(): %v", err)
	}
	if err := activityEq(got[1], activities[1]); err != nil {
		t.Fatalf("GetAllActivities(): %v", err)
	}
	if err := activityEq(got[len(got)-1], activities[len(activities)-1]); err != nil {
		t.Fatalf("GetAllActivities(): %v", err)
	}

	got, err = db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
		OrganizationId: organizationB,
	})
	if err != nil {
		t.Fatalf("GetAllActivities(): got error; want nil")
	}
	if len(got) != 0 {
		t.Fatalf("GetAllActivities(): got %d activities; want %d activities", len(got), 0)
	}

	db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
		Id:             activities[0].Id,
		OrganizationId: organizationA,
	})
	got, err = db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
		OrganizationId: organizationA,
	})
	if err != nil {
		t.Fatalf("GetAllActivities(): got error; want nil")
	}
	if len(got) != NUM_ACTIVITIES_CREATED-1 {
		t.Fatalf("GetAllActivities(): got %d activities; want %d activities", len(got), NUM_ACTIVITIES_CREATED)
	}
	if err := activityEq(got[0], activities[1]); err != nil {
		t.Fatalf("GetAllActivities(): %v", err)
	}
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
