package storage

import (
	"context"

	"stockinos.com/api/models"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg DoesUserExistsParams) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, arg GetUserByPhoneNumberParams) (*models.User, error)
	UpdateUserName(ctx context.Context, arg UpdateUserNameParams) (*models.User, error)

	// OTP
	CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.OTP, error)
	DesactivateOTP(ctx context.Context, arg DesactivateOTPParams) (*models.OTP, error)
	DesactivateAllOTPFromPhoneNumber(ctx context.Context, arg DesactivateAllOTPFromPhoneNumberParams) error
	CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.OTP, error)
	GetActivateOTP(ctx context.Context, arg GetActivateOTPParams) (*models.OTP, error)

	// Company
	CreateCompany(ctx context.Context, arg CreateCompanyParams) (*models.Company, error)
	GetAllCompanies(ctx context.Context, arg GetAllCompaniesParams) ([]*models.Company, error)
	GetCompany(ctx context.Context, arg GetCompanyParams) (*models.Company, error)
	UpdateCompany(ctx context.Context, arg UpdateCompanyParams) (*models.Company, error)
	DeleteCompany(ctx context.Context, arg DeleteCompanyParams) error

	// Monitoring
	// Activity
	CreateActivity(ctx context.Context, arg CreateActivityParams) (*models.Activity, error)
	GetActivity(ctx context.Context, arg GetActivityParams) (*models.Activity, error)
	DeleteActivity(ctx context.Context, arg DeleteActivityParams) error
	GetAllActivities(ctx context.Context, arg GetAllActivitiesParams) ([]*models.Activity, error)
	UpdateSetInActivity(ctx context.Context, arg UpdateSetInActivityParams) (*models.Activity, error)
	UpdateAddToActivity(ctx context.Context, arg UpdateAddToActivityParams) (*models.Activity, error)
	UpdateRemoveFromActivity(ctx context.Context, arg UpdateRemoveFromActivityParams) (*models.Activity, error)
	// Data
	CreateData(ctx context.Context, arg CreateDataParams) (*models.Data, error)
	GetData(ctx context.Context, arg GetDataParams) (*models.Data, error)
	GetAllData(ctx context.Context, arg GetAllDataParams) ([]*models.Data, error)
	DeleteData(ctx context.Context, arg DeleteDataParams) error
	UpdateSetInData(ctx context.Context, arg UpdateSetInDataParams) (*models.Data, error)
	UpdateAddToData(ctx context.Context, arg UpdateAddToDataParams) (*models.Data, error)
	UpdateRemoveFromData(ctx context.Context, arg UpdateRemoveFromDataParams) (*models.Data, error)
}

type QuerierTx interface {
	// OTP
	CheckOTPTx(ctx context.Context, arg CheckOTPParams) (*models.OTP, error)
	CreateOTPx(ctx context.Context, arg CreateOTPParams) (*models.OTP, error)
}

var _ Querier = (*Queries)(nil)
