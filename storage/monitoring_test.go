package storage_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		CreatedBy: primitive.NewObjectID(),
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
	if activity.CreatedBy != arg.CreatedBy {
		t.Fatalf("CreateActivity() createdBy - got: %v; want: %v", activity.CreatedBy, arg.CreatedBy)
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
	if got.CreatedBy != arg.CreatedBy {
		t.Fatalf("CreateActivity() createdBy - got: %v; want: %v", got.CreatedBy, arg.CreatedBy)
	}
}

func TestDeleteActivity(t *testing.T) {

	db, cleanup := integrationtest.CreateDatabase()
	defer cleanup()

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

	_, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id: activity.Id,
	})
	if err != nil {
		t.Fatalf("GetActivity() - err = %v; want nil", err)
	}

	err = db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
		Id: activity.Id,
	})
	if err != nil {
		t.Fatalf("DeleteActivity() - err = %v; want nil", err)
	}

	activities, err := db.Storage.GetAllActivitiesFromUser(context.Background(), storage.GetAllActivitiesFromUserParams{
		CreatedBy: arg.CreatedBy,
	})
	if err != nil {
		t.Fatalf("GetAllActivitiesFromUser() - err = %v; want nil", err)
	}
	if len(activities) != 0 {
		t.Fatalf("GetAllActivitiesFromUser() size - got: %d; want = 0", len(activities))
	}
}
