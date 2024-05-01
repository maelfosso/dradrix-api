package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

func (store *MongoStorage) CheckOTPTx(ctx context.Context, arg CheckOTPParams) (*models.OTP, error) {
	result, err := store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		otp, err := store.CheckOTP(ctx, CheckOTPParams{
			PhoneNumber: arg.PhoneNumber,
			UserOTP:     arg.UserOTP,
		})
		if err != nil {
			log.Println("error when checking the otp: ", err)
			return nil, fmt.Errorf("ERR_COTP_102_%s", err)
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
