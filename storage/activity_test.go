package storage_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

func TestActivity(t *testing.T) {
	integrationtest.SkipifShort(t)

	db, disconnect := integrationtest.CreateDatabase()
	defer disconnect()

	tests := map[string]func(*testing.T, *storage.Database){
		"CreateActivity": testCreateActivity,
		"DeleteActivity": testDeleteActivity,
		"UpdateActivity": testUpdateActivity,
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
		Fields: []models.ActivityFields{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text"},
		},

		CompanyId: primitive.NewObjectID(),
		CreatedBy: primitive.NewObjectID(),
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

	got, err := db.Storage.GetActivity(context.TODO(), storage.GetActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("GetActivityFromCompany(): got error %+v; want nit", err.Error())
	}
	if err := activityEq(got, activity); err != nil {
		t.Fatalf("GetActivityFromCompany(): %v", err.Error())
	}
}

func testDeleteActivity(t *testing.T, db *storage.Database) {
	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityFields{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text"},
		},

		CompanyId: primitive.NewObjectID(),
		CreatedBy: primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity(): err = %v; want nil", err)
	}

	_, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("GetActivityFromCompany(): got error = %v; want nil", err)
	}

	err = db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("DeleteActivityFromCompany(): got err = %v; want nil", err)
	}

	activity, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("GetActivityFromCompany(): got error = %v; want nil", err)
	}
	if activity != nil {
		t.Fatalf("GetActivityFromCompany(): got %v; want nil", activity)
	}

	activities, err := db.Storage.GetAllActivities(context.Background(), storage.GetAllActivitiesParams{
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("GetAllActivitiesFromCompany(): got err = %v; want nil", err)
	}
	if len(activities) != 0 {
		t.Fatalf("GetAllActivitiesFromCompany(): got %d number of activities; want = 0", len(activities))
	}
}

func testUpdateActivity(t *testing.T, db *storage.Database) {

	beforeCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() Before; err = %v; want nil", err)
	}

	arg := storage.CreateActivityParams{
		Name:        "a1",
		Description: "Activity 1",
		Fields: []models.ActivityFields{
			{Name: "f1", Description: "Description 1", Type: "number"},
			{Name: "f2", Description: "Description 2", Type: "text", Id: true},
		},

		CompanyId: primitive.NewObjectID(),
		CreatedBy: primitive.NewObjectID(),
	}
	activity, err := db.Storage.CreateActivity(context.Background(), arg)
	if err != nil {
		t.Fatalf("CreateActivity() err = %v; want nil", err)
	}
	if activity.Id.IsZero() {
		t.Fatalf("CreateActivity(): Id is nil; want non nil")
	}

	argForUpdate := storage.UpdateActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,

		Field: "name",
		Value: "a2",
	}
	updated, err := db.Storage.UpdateActivity(context.TODO(), argForUpdate)
	if err != nil {
		t.Fatalf("UpdateActivity(): got error %v; want nil", err)
	}
	if updated.Name != argForUpdate.Value {
		t.Fatalf(
			"UpdateActivity(): updated %s value - got %s; want %s",
			argForUpdate.Field, updated.Name, argForUpdate.Value,
		)
	}

	argForUpdate = storage.UpdateActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,

		Field: "fields.1.code",
		Value: "f1",
	}
	updated, err = db.Storage.UpdateActivity(context.TODO(), argForUpdate)
	if err != nil {
		t.Fatalf("UpdateActivity(): got error %v; want nil", err)
	}
	if updated.Fields[1].Code != argForUpdate.Value {
		t.Fatalf(
			"UpdateActivity(): updated %s value - got %s; want %s",
			argForUpdate.Field, updated.Name, argForUpdate.Value,
		)
	}

	argForUpdate = storage.UpdateActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,

		Field: "fields.1.id",
		Value: false,
	}
	updated, err = db.Storage.UpdateActivity(context.TODO(), argForUpdate)
	if err != nil {
		t.Fatalf("UpdateActivity(): got error %v; want nil", err)
	}
	if updated.Fields[0].Id != argForUpdate.Value {
		t.Fatalf(
			"UpdateActivity(): updated %s value - got %s; want %s",
			argForUpdate.Field, updated.Name, argForUpdate.Value,
		)
	}

	afterCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
	if err != nil {
		t.Fatalf("CountDocuments() After; err = %v; want nil", err)
	}
	if afterCount-beforeCount != 1 {
		t.Fatalf("AfterCount - BeforeCount = %d; want = %d", afterCount-beforeCount, 1)
	}

	got, err := db.Storage.GetActivity(context.TODO(), storage.GetActivityParams{
		Id:        activity.Id,
		CompanyId: arg.CompanyId,
	})
	if err != nil {
		t.Fatalf("GetActivityFromCompany(): got error %+v; want nit", err.Error())
	}
	if err := activityEq(got, updated); err != nil {
		t.Fatalf("GetActivityFromCompany(): %v", err.Error())
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

// func TestDeleteActivity(t *testing.T) {

// 	db, cleanup := integrationtest.CreateDatabase()
// 	defer cleanup()

// 	arg := storage.CreateActivityParams{
// 		Name:        "a1",
// 		Description: "Activity 1",
// 		Fields: []models.ActivityFields{
// 			{Name: "f1", Description: "Description 1", Type: "number"},
// 			{Name: "f2", Description: "Description 2", Type: "text"},
// 		},
// 	}
// 	activity, err := db.Storage.CreateActivity(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("CreateActivity() err = %v; want nil", err)
// 	}

// 	_, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
// 		Id: activity.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetActivity() - err = %v; want nil", err)
// 	}

// 	err = db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
// 		Id: activity.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("DeleteActivity() - err = %v; want nil", err)
// 	}

// 	activities, err := db.Storage.GetAllActivitiesFromUser(context.Background(), storage.GetAllActivitiesFromUserParams{
// 		CreatedBy: arg.CreatedBy,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivitiesFromUser() - err = %v; want nil", err)
// 	}
// 	if len(activities) != 0 {
// 		t.Fatalf("GetAllActivitiesFromUser() size - got: %d; want = 0", len(activities))
// 	}
// }

// func TestCreateData(t *testing.T) {
// 	integrationtest.SkipifShort(t)

// 	db, cleanup := integrationtest.CreateDatabase()
// 	defer cleanup()

// 	beforeCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
// 	if err != nil {
// 		t.Fatalf("CountDocuments() Before; err = %v; want nil", err)
// 	}

// 	arg := storage.CreateActivityParams{
// 		Name:        "a1",
// 		Description: "Activity 1",
// 		Fields: []models.ActivityFields{
// 			{Name: "f1", Description: "Description 1", Type: "number"},
// 			{Name: "f2", Description: "Description 2", Type: "text"},
// 		},

// 		CreatedBy: primitive.NewObjectID(),
// 	}
// 	activity, err := db.Storage.CreateActivity(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("CreateActivity() err = %v; want nil", err)
// 	}
// 	if activity.Id.IsZero() {
// 		t.Fatalf("CreateActivity(); Id is nil; want non nil")
// 	}
// 	if len(activity.Name) == 0 && len(activity.Description) == 0 && len(activity.Fields) <= 0 {
// 		t.Fatalf("CreateActivity(): Properties are not okay; got = (%s, %s, %d); want = (%s, %s, %d)",
// 			activity.Name, activity.Description, len(activity.Fields),
// 			arg.Name, arg.Description, len(arg.Fields),
// 		)
// 	}
// 	if activity.CreatedAt.IsZero() || activity.UpdatedAt.IsZero() {
// 		t.Fatalf("CreateActivity() date - got empty date; want date with values")
// 	}
// 	if activity.CreatedBy != arg.CreatedBy {
// 		t.Fatalf("CreateActivity() createdBy - got: %v; want: %v", activity.CreatedBy, arg.CreatedBy)
// 	}

// 	afterCount, err := db.GetCollection("activities").CountDocuments(context.Background(), bson.D{}, nil)
// 	if err != nil {
// 		t.Fatalf("CountDocuments() After; err = %v; want nil", err)
// 	}
// 	if afterCount-beforeCount != 1 {
// 		t.Fatalf("AfterCount - BeforeCount = %d; want = %d", afterCount-beforeCount, 1)
// 	}

// 	got, err := db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
// 		Id: activity.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetActivity() - err = %v; want nil", err)
// 	}
// 	if got.Name != arg.Name || got.Description != arg.Description || len(got.Fields) != len(arg.Fields) {
// 		t.Fatalf("GetActivity - got: (%s, %s, %d); want: (%s, %s, %d)",
// 			got.Name, got.Description, len(got.Fields),
// 			arg.Name, arg.Description, len(arg.Fields),
// 		)
// 	}
// 	if got.CreatedBy != arg.CreatedBy {
// 		t.Fatalf("CreateActivity() createdBy - got: %v; want: %v", got.CreatedBy, arg.CreatedBy)
// 	}
// }

// func TestDeleteData(t *testing.T) {

// 	db, cleanup := integrationtest.CreateDatabase()
// 	defer cleanup()

// 	arg := storage.CreateActivityParams{
// 		Name:        "a1",
// 		Description: "Activity 1",
// 		Fields: []models.ActivityFields{
// 			{Name: "f1", Description: "Description 1", Type: "number"},
// 			{Name: "f2", Description: "Description 2", Type: "text"},
// 		},
// 	}
// 	activity, err := db.Storage.CreateActivity(context.Background(), arg)
// 	if err != nil {
// 		t.Fatalf("CreateActivity() err = %v; want nil", err)
// 	}

// 	_, err = db.Storage.GetActivity(context.Background(), storage.GetActivityParams{
// 		Id: activity.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetActivity() - err = %v; want nil", err)
// 	}

// 	err = db.Storage.DeleteActivity(context.Background(), storage.DeleteActivityParams{
// 		Id: activity.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("DeleteActivity() - err = %v; want nil", err)
// 	}

// 	activities, err := db.Storage.GetAllActivitiesFromUser(context.Background(), storage.GetAllActivitiesFromUserParams{
// 		CreatedBy: arg.CreatedBy,
// 	})
// 	if err != nil {
// 		t.Fatalf("GetAllActivitiesFromUser() - err = %v; want nil", err)
// 	}
// 	if len(activities) != 0 {
// 		t.Fatalf("GetAllActivitiesFromUser() size - got: %d; want = 0", len(activities))
// 	}
// }
