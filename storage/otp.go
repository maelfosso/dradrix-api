package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"stockinos.com/api/models"
)

type GetActivateOTPParams struct {
	PhoneNumber string
}

func (q *Queries) GetActivateOTP(ctx context.Context, arg GetActivateOTPParams) (*models.OTP, error) {
	var otp models.OTP
	err := q.otpsCollection.FindOne(
		ctx,
		bson.D{{Key: "phone_number", Value: arg.PhoneNumber}, {Key: "active", Value: true}},
	).Decode(&otp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &otp, nil
}

type CreateOTPParams struct {
	WaMessageId string
	PinCode     string
	PhoneNumber string
}

func (q *Queries) CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.OTP, error) {
	var otp models.OTP = models.OTP{
		Id:          primitive.NewObjectID(),
		PhoneNumber: arg.PhoneNumber,
		WaMessageId: arg.WaMessageId,
		PinCode:     arg.PinCode,
		Active:      true,
	}

	_, err := q.otpsCollection.InsertOne(ctx, otp)
	if err != nil {
		return nil, err
	} else {
		return &otp, nil
	}
}

type DesactivateOTPParams struct {
	Id primitive.ObjectID
}

func (q *Queries) DesactivateOTP(ctx context.Context, arg DesactivateOTPParams) (*models.OTP, error) {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}
	after := options.After

	var otp models.OTP
	err := q.otpsCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&otp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &otp, nil
}

type DesactivateAllOTPFromPhoneNumberParams struct {
	PhoneNumber string
}

func (q *Queries) DesactivateAllOTPFromPhoneNumber(ctx context.Context, arg DesactivateAllOTPFromPhoneNumberParams) error {
	filter := bson.M{
		"phone_number": arg.PhoneNumber,
	}
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}

	_, err := q.otpsCollection.UpdateMany(
		ctx,
		filter,
		update,
	)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

type CheckOTPParams struct {
	PhoneNumber string
	UserOTP     string
}

func (q *Queries) CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.OTP, error) {
	var otp models.OTP

	filter := bson.M{
		"phone_number": arg.PhoneNumber,
		"pin_code":     arg.UserOTP,
		"active":       true,
	}
	err := q.otpsCollection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &otp, nil
}
