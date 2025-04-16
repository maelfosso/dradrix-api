package storage_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

func TestCreateOTP(t *testing.T) {
	integrationtest.SkipifShort(t)

	t.Run("create otp", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		otp, err := db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "0000-0000-0",
		})
		if err != nil {
			t.Fatalf("CreateOTP() error %v", err)
		}
		if otp == nil {
			t.Fatalf("CreateOTP() otp is nil; want not nil")
		}

		err = db.GetCollection("otps").FindOne(context.Background(), bson.M{
			"phone_number": "0000-0000-0",
			"pin_code":     "0000",
			"active":       true,
		}).Decode(&otp)
		if err != nil {
			t.Fatalf("CreateOTP() - Check result; err = %v; want = nil", err)
		}
		if otp == nil {
			t.Fatalf("CreateOTP - Check result; otp = nil; want = 0000-0000-0")
		}
	})
}

func TestGetActivateOTP(t *testing.T) {
	integrationtest.SkipifShort(t)

	t.Run("Get activate otp", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		otpCreated, _ := db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "1111-0000-0",
		})
		otp, err := db.Storage.GetActivateOTP(context.Background(), storage.GetActivateOTPParams{
			PhoneNumber: "1111-0000-0",
		})
		if err != nil {
			t.Fatalf("GetActivateOTP() err = %v; want = nil", err)
		}
		if otp.PhoneNumber != otpCreated.PhoneNumber {
			t.Fatalf("GetActivateOTP() Fetched phone number = %s; want = %s", otp.PhoneNumber, otp.PhoneNumber)
		}
		if otp.Active != true {
			t.Fatalf("GetActivateOTP() Fetched OTP Active = %t; want = true", otp.Active)
		}
	})
}

func TestDesactivateOTP(t *testing.T) {
	t.Run("Desactivate otp", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		otp, _ := db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "2222-2222-2",
		})
		_, err := db.Storage.DesactivateOTP(context.Background(), storage.DesactivateOTPParams{
			Id: otp.Id,
		})
		if err != nil {
			t.Fatalf("DesactivateOTP() err=%v, want=nil", err)
		}

		var otpFetched *models.OTP
		err = db.GetCollection("otps").FindOne(context.Background(), bson.M{
			"phone_number": "2222-2222-2",
			"pin_code":     "0000",
			"active":       true,
		}).Decode(otpFetched)
		if err != nil && err != mongo.ErrNoDocuments {
			t.Fatalf("DesactivateOTP() otp still active err=%v, want=nil", err)
		}
		if otpFetched != nil {
			t.Fatalf("DesactivateOTP() otp Id=%s; want=nil", otpFetched.Id.String())
		}
	})
}

func TestCheckOTP(t *testing.T) {
	t.Run("CheckOTP()", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		_, _ = db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "0000-0000-0",
		})

		otp, err := db.Storage.CheckOTP(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "0000-0000-0",
			UserOTP:     "0000",
		})
		if err != nil {
			t.Fatalf("CheckOTP err=%v, want=nil", err)
		}
		if otp.PhoneNumber != "0000-0000-0" || otp.PinCode != "0000" || !otp.Active {
			t.Fatalf("CheckOTP wrong fetched otp=%s/%s/%t; want=0000-0000-0/0000/true", otp.PhoneNumber, otp.PinCode, otp.Active)
		}
	})

	t.Run("CheckOTP() with wrong OTP/PhoneNumber", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		_, _ = db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "0000-0000-0",
		})

		otp, err := db.Storage.CheckOTP(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "0000-0000-0",
			UserOTP:     "1111",
		})
		if err != nil {
			t.Fatalf("CheckOTP err=%v, want=nil", err)
		}
		if otp != nil {
			t.Fatalf("CheckOTP wrong otp: one otp fetched=%s/%s; want=nil", otp.PhoneNumber, otp.PinCode)
		}

		otp, err = db.Storage.CheckOTP(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "1111-0000-0",
			UserOTP:     "0000",
		})
		if err != nil {
			t.Fatalf("CheckOTP err=%v, want=nil", err)
		}
		if otp != nil {
			t.Fatalf("CheckOTP wrong phone number: one otp fetched=%s/%s; want=nil", otp.PhoneNumber, otp.PinCode)
		}
	})
}
