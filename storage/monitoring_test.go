package storage_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

func TestCreateActivity(t *testing.T) {
	integrationtest.SkipifShort(t)

	db, cleanup := integrationtest.CreateDatabase()
	defer cleanup()

	beforeCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() Before; err = %v; want nil", err)
	}

	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityFields{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text"},
		},
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity() err = %v; want nil", err)
	}
	if activity.Id.IsZero() {
		t.Fatalf("CreateActivity(); Id is nil; want non nil")
	}
	if len(activity.Name) == 0 && len(activity.Description) == 0 && len(activity.Fields) <= 0 {
		t.Fatalf("CreateActivity(): Properties are not okay; got = (%s, %s, %d); want = (%s, %s, %d)",
			activity.Name, activity.Description, len(activity.Fields),
			arg.Name, arg.Description, len(arg.Fields),
		)
	}
	if activity.CreatedAt.IsZero() || activity.UpdatedAt.IsZero() {
		t.Fatalf("CreateActivity() date - got empty date; want date with values")
	}

	afterCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() After; err = %v; want nil", err)
	}
	if afterCount-beforeCount != 1 {
		t.Fatalf("AfterCount - BeforeCount = %d; want = %d", afterCount-beforeCount, 1)
	}

	got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id: activity.Id,
	})
	if err != nil {
		t.Fatalf("GetActivity() - err = %v; want nil", err)
	}
	if got.Name != arg.Name || got.Description != arg.Description || len(got.Fields) != len(arg.Fields) {
		t.Fatalf("GetActivity - got: (%s, %s, %d); want: (%s, %s, %d)",
			got.Name, got.Description, len(got.Fields),
			arg.Name, arg.Description, len(arg.Fields),
		)
	}

}
