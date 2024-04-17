package storage

import (
	"context"

	"stockinos.com/api/models"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg DoesUserExistsParams) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, arg GetUserByPhoneNumberParams) (*models.User, error)

	// OTP
	CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.OTP, error)
	DesactivateOTP(ctx context.Context, arg DesactivateOTPParams) error
	CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.OTP, error)
	GetActivateOTP(ctx context.Context, arg GetActivateOTPParams) (*models.OTP, error)
}

type QuerierTx interface {
}

var _ Querier = (*Queries)(nil)
