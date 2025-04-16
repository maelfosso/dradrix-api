package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

func (store *MongoStorage) CreateOTPx(ctx context.Context, arg CreateOTPParams) (*models.OTP, error) {
	result, err := store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := store.DesactivateAllOTPFromPhoneNumber(ctx, DesactivateAllOTPFromPhoneNumberParams{
			PhoneNumber: arg.PhoneNumber,
		})
		if err != nil {
			log.Println("error when desactivating all the otp ", err)
			return nil, fmt.Errorf("%v", err)
		}

		waMessageId := "xxx-yyy-zzz"
		otp, err := store.CreateOTP(ctx, CreateOTPParams{
			WaMessageId: waMessageId,
			PhoneNumber: arg.PhoneNumber,
			PinCode:     arg.PinCode,
		})
		if err != nil {
			log.Println("error when saving the OTP: ", err)
			return nil, fmt.Errorf("%v", err)
		}

		return otp, nil
	})

	if otp, ok := result.(*models.OTP); ok {
		return otp, err
	} else {
		return nil, err
	}
}

func (store *MongoStorage) CheckOTPTx(ctx context.Context, arg CheckOTPParams) (*models.OTP, error) {
	result, err := store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		otp, err := store.CheckOTP(ctx, CheckOTPParams{
			PhoneNumber: arg.PhoneNumber,
			UserOTP:     arg.UserOTP,
		})
		if err != nil {
			log.Println("error when checking the otp: ", err)
			return nil, fmt.Errorf("ERR_COTP_102_%v", err)
		}
		if otp == nil {
			log.Println("error when checking the otp - no corresponding otp found: ", err)
			return nil, fmt.Errorf("ERR_COTP_102_%s", err)
		}

		otp, err = store.DesactivateOTP(ctx, DesactivateOTPParams{
			Id: otp.Id,
		})
		if err != nil {
			return nil, fmt.Errorf("ERR_", err)
		}

		return otp, nil
	})

	if otp, ok := result.(*models.OTP); ok {
		return otp, err
	} else {
		return nil, err
	}
}
