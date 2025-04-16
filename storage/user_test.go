package storage_test

import (
	"context"
	"fmt"
	"testing"

	"stockinos.com/api/integrationtest"
	"stockinos.com/api/storage"
)

func TestCreateUser(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		user, err := db.Storage.CreateUser(context.Background(), storage.CreateUserParams{
			PhoneNumber: "0000-0000",
			FirstName:   "Doe",
			LastName:    "John",
		})
		if err != nil {
			t.Fatalf("CreateUser() - err=%v; want=nil", err)
		}
		if user == nil {
			t.Fatalf("CreateUser() user should not be nil")
		}

		user, err = db.Storage.GetUserByPhoneNumber(context.Background(), storage.GetUserByPhoneNumberParams{
			PhoneNumber: "0000-0000",
		})
		if err != nil {
			t.Fatalf("GetUserByPhoneNumber() - err=%v; want=nil", err)
		}
		if user == nil {
			t.Fatalf("GetUserByPhoneNumber() user should not be nil")
		}
		if user.PhoneNumber != "0000-0000" || user.FirstName != "Doe" || user.LastName != "John" {
			t.Fatalf("GetUserByPhoneNumber() - wrong fetched user; got=%s/%s; want=0000-0000/John Doe", user.PhoneNumber, fmt.Sprintf("%s %s", user.LastName, user.FirstName))
		}

		user, err = db.Storage.DoesUserExists(context.Background(), storage.DoesUserExistsParams{
			PhoneNumber: "0000-0000",
		})
		if err != nil {
			t.Fatalf("DoesUserExists() - err=%v; want=nil", err)
		}
		if user == nil {
			t.Fatalf("DoesUserExists() user should not be nil")
		}
		if user.PhoneNumber != "0000-0000" || user.FirstName != "Doe" || user.LastName != "John" {
			t.Fatalf("DoesUserExists() - wrong fetched user; got=%s/%s; want=0000-0000/John Doe", user.PhoneNumber, fmt.Sprintf("%s %s", user.LastName, user.FirstName))
		}
	})
}

func TestGetUserByPhoneNumber(t *testing.T) {
	t.Run("GetUserByPhoneNumber()", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		db.Storage.CreateUser(context.Background(), storage.CreateUserParams{
			PhoneNumber: "0000-0000",
			FirstName:   "Doe",
			LastName:    "John",
		})

		user, err := db.Storage.GetUserByPhoneNumber(context.Background(), storage.GetUserByPhoneNumberParams{
			PhoneNumber: "1111-0000",
		})
		if err != nil {
			t.Fatalf("GetUserByPhoneNumber() - err=%v; want=nil", err)
		}
		if user != nil {
			t.Fatalf("GetUserByPhoneNumber() user should be nil")
		}
	})
}

func TestDoesUserExists(t *testing.T) {
	t.Run("DoesUserExists()", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		db.Storage.CreateUser(context.Background(), storage.CreateUserParams{
			PhoneNumber: "0000-0000",
			FirstName:   "Doe",
			LastName:    "John",
		})

		user, err := db.Storage.DoesUserExists(context.Background(), storage.DoesUserExistsParams{
			PhoneNumber: "1111-0000",
		})
		if err != nil {
			t.Fatalf("DoesUserExists() - err=%v; want=nil", err)
		}
		if user != nil {
			t.Fatalf("DoesUserExists() user should be nil")
		}
	})
}
