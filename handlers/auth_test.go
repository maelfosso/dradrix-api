package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/handlers"
	"stockinos.com/api/helpertest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type getOTPMock struct{}

var userId = primitive.NewObjectID()
var otpId = primitive.NewObjectID()
var pinCode = "0000"
var users []models.User
var otps []models.OTP

func (s *getOTPMock) CreateUser(ctx context.Context, arg storage.CreateUserParams) (*models.User, error) {
	users = append(users, models.User{
		Id:          userId,
		PhoneNumber: "695165033",
	})

	return &users[len(users)-1], nil
}

func (s *getOTPMock) DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
	var user *models.User = nil
	for _, _user := range users {
		if _user.PhoneNumber == arg.PhoneNumber {
			user = &_user

			break
		}
	}
	return user, nil
}

func (s *getOTPMock) CreateOTP(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
	otp := models.OTP{
		Id:          otpId,
		WaMessageId: "xxx-yyy-zzz",
		PhoneNumber: "695165033",
		PinCode:     pinCode,
		Active:      true,
	}
	otps = append(otps, otp)

	return &otps[len(otps)-1], nil
}

func Init() {
	users = make([]models.User, 0)
	otps = make([]models.OTP, 0)
}

func TestCreateOTP(t *testing.T) {
	mux := chi.NewMux()
	svc := &getOTPMock{}
	handlers.CreateOTP(mux, svc)

	t.Run("return 200", func(t *testing.T) {
		Init()

		code, _, _ := helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		if code != http.StatusOK {
			t.Fatalf("CreateOTP() status code = %d; want = %d", code, http.StatusOK)
		}
		if users[len(users)-1].PhoneNumber != "695165033" {
			t.Fatalf("CreateOTP() last user phone number = %s; want = %s", otps[len(otps)-1].PhoneNumber, "695165033")
		}
		if otps[len(otps)-1].PhoneNumber != "695165033" {
			t.Fatalf("CreateOTP() created otp phone number = %s; want = %s", otps[len(otps)-1].PhoneNumber, "695165033")
		}
		if otps[len(otps)-1].Active != true {
			t.Fatalf("CreateOTP() create otp active status = %v; want = %v", otps[len(otps)-1].Active, true)
		}
	})

	t.Run("return 200 but with no more users if phone number already exists", func(t *testing.T) {
		Init()
		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})
		code, _, _ := helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})
		if code != http.StatusOK {
			t.Fatalf("CreateOTP() status code = %d; want = %d", code, http.StatusOK)
		}
		if len(users) == 2 {
			t.Fatalf("CreateOTP() number of users created = %d; want = %d", len(users), 1)
		}
	})

}
