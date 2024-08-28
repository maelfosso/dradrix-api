package storage

import (
	"context"

	"stockinos.com/api/models"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error)
	DoesUserExists(ctx context.Context, arg DoesUserExistsParams) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, arg GetUserByPhoneNumberParams) (*models.User, error)
	UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (*models.User, error)
	UpdateUserPreferences(ctx context.Context, arg UpdateUserPreferencesParams) (*models.User, error)

	// OTP
	CheckOTP(ctx context.Context, arg CheckOTPParams) (*models.OTP, error)
	DesactivateOTP(ctx context.Context, arg DesactivateOTPParams) (*models.OTP, error)
	DesactivateAllOTPFromPhoneNumber(ctx context.Context, arg DesactivateAllOTPFromPhoneNumberParams) error
	CreateOTP(ctx context.Context, arg CreateOTPParams) (*models.OTP, error)
	GetActivateOTP(ctx context.Context, arg GetActivateOTPParams) (*models.OTP, error)

	// Organization
	CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error)
	GetAllCompanies(ctx context.Context, arg GetAllCompaniesParams) ([]*models.Organization, error)
	GetOrganization(ctx context.Context, arg GetOrganizationParams) (*models.Organization, error)
	UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, arg DeleteOrganizationParams) error

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
	UpdateData(ctx context.Context, arg UpdateDataParams) (*models.Data, error)
	GetData(ctx context.Context, arg GetDataParams) (*models.Data, error)
	GetDataFromValues(ctx context.Context, arg GetDataFromValuesParams) (*models.Data, error)
	GetAllData(ctx context.Context, arg GetAllDataParams) ([]*models.Data, error)
	DeleteData(ctx context.Context, arg DeleteDataParams) error
	UpdateSetInData(ctx context.Context, arg UpdateSetInDataParams) (*models.Data, error)
	UpdateAddToData(ctx context.Context, arg UpdateAddToDataParams) (*models.Data, error)
	UpdateRemoveFromData(ctx context.Context, arg UpdateRemoveFromDataParams) (*models.Data, error)
	AddUploadedFile(ctx context.Context, arg AddUploadedFileParams) (*models.UploadedFile, error)
	GetAllUploadedFiles(ctx context.Context, arg GetAllUploadedFilesParams) ([]*models.UploadedFile, error)
	RemoveUploadedFile(ctx context.Context, arg RemoveUploadedFileParams) error
	RemoveAllUploadedFile(ctx context.Context, arg RemoveUploadedFileParams) error
}

type QuerierTx interface {
	// OTP
	CheckOTPTx(ctx context.Context, arg CheckOTPParams) (*models.OTP, error)
	CreateOTPx(ctx context.Context, arg CreateOTPParams) (*models.OTP, error)
}

var _ Querier = (*Queries)(nil)
