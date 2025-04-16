package storage_test

import (
	"context"
	"testing"

	"stockinos.com/api/integrationtest"
	"stockinos.com/api/storage"
)

func TestCheckOTPTex(t *testing.T) {
	t.Run("CheckOTPTx()", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		_, _ = db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "0000-0000-0",
		})

		otp, err := db.Storage.CheckOTPTx(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "0000-0000-0",
			UserOTP:     "0000",
		})
		if err != nil {
			t.Fatalf("CheckOTPTx err=%v, want=nil", err)
		}
		if otp.PhoneNumber != "0000-0000-0" && otp.PinCode != "0000" {
			t.Fatalf("CheckOTP wrong fetched otp=%s/%s; want=0000-0000-0/0000", otp.PhoneNumber, otp.PinCode)
		}
		if otp.Active {
			t.Fatalf("CheckOTP fetched otp is still active; active=true; want=false")
		}
	})
	t.Run("CheckOTPTx() with bad OTP/PhoneNumber", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		_, _ = db.Storage.CreateOTP(context.Background(), storage.CreateOTPParams{
			WaMessageId: "xxxx",
			PinCode:     "0000",
			PhoneNumber: "0000-0000-0",
		})

		otp, err := db.Storage.CheckOTPTx(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "0000-0000-0",
			UserOTP:     "1111",
		})
		if err == nil {
			t.Fatal("CheckOTP err=nil, want=error")
		}
		if otp != nil {
			t.Fatalf("CheckOTP wrong otp: one otp fetched=%s/%s; want=nil", otp.PhoneNumber, otp.PinCode)
		}

		otp, err = db.Storage.CheckOTPTx(context.Background(), storage.CheckOTPParams{
			PhoneNumber: "1111-0000-0",
			UserOTP:     "0000",
		})
		if err == nil {
			t.Fatal("CheckOTP err=nil, want=error")
		}
		if otp != nil {
			t.Fatalf("CheckOTP wrong phone number: one otp fetched=%s/%s; want=nil", otp.PhoneNumber, otp.PinCode)
		}
	})
}
