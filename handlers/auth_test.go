package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/handlers"
	"stockinos.com/api/helpertest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

type getOTPMock struct{}

var userId = primitive.NewObjectID()
var otpId = primitive.NewObjectID()
var pinCode = utils.GenerateOTP(time.Now())
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

func (s *getOTPMock) CreateOTPx(ctx context.Context, arg storage.CreateOTPParams) (*models.OTP, error) {
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

		userPhoneNumber := "695165033"
		code, _, responseData := helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: userPhoneNumber,
			Language:    "fr",
		})

		if code != http.StatusOK {
			t.Fatalf("CreateOTP() status code = %d; want = %d", code, http.StatusOK)
		}

		var responsePhoneNumber string
		err := json.Unmarshal([]byte(responseData), &responsePhoneNumber)
		if err != nil || responsePhoneNumber != userPhoneNumber {
			t.Fatalf("CreateOTP() response request = %s; want = %s", responsePhoneNumber, userPhoneNumber)
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

	t.Run("return 200 but always only one active otp", func(t *testing.T) {
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
		if len(otps) < 2 {
			t.Fatalf("CreateOTP() number of otps created = %d; want = %d", len(otps), 2)
		}
		nActiveOTP := 0
		for i := range otps {
			otp := otps[i]
			if otp.Active == true {
				nActiveOTP += 1
			}
		}
		if nActiveOTP != 1 {
			t.Fatalf("CreateOTP number of active otps = %d; want = %d", nActiveOTP, 1)
		}
	})

}

type checkOTPMock struct{}

func (s *checkOTPMock) DoesUserExists(ctx context.Context, arg storage.DoesUserExistsParams) (*models.User, error) {
	var user *models.User = nil
	for _, _user := range users {
		if _user.PhoneNumber == arg.PhoneNumber {
			user = &_user

			break
		}
	}
	return user, nil
}

func (s *checkOTPMock) CheckOTPTx(ctx context.Context, arg storage.CheckOTPParams) (*models.OTP, error) {
	var otp *models.OTP = nil
	var ix int
	for _i, _otp := range otps {
		if _otp.PhoneNumber == arg.PhoneNumber && _otp.PinCode == arg.UserOTP && _otp.Active == true {
			otp = &_otp
			ix = _i

			break
		}
	}
	if otp != nil {
		otps[ix].Active = false
	}
	return otp, nil
}

func TestCheckOTP(t *testing.T) {
	mux := chi.NewMux()
	checkSVC := &checkOTPMock{}
	getSVC := &getOTPMock{}
	handlers.CheckOTP(mux, checkSVC)
	handlers.CreateOTP(mux, getSVC)

	t.Run("return 200", func(t *testing.T) {
		Init()

		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		code, _, _ := helpertest.MakePostRequest(mux, "/otp/check", helpertest.CreateFormHeader(), handlers.CheckOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
			PinCode:     otps[len(otps)-1].PinCode,
		})
		if code != http.StatusOK {
			t.Fatalf("CheckOTP() status code = %d; want = %d", code, http.StatusOK)
		}
	})

	t.Run("return 400 if the user phone number doesn't exists", func(t *testing.T) {
		Init()

		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		code, _, _ := helpertest.MakePostRequest(mux, "/otp/check", helpertest.CreateFormHeader(), handlers.CheckOTPRequest{
			PhoneNumber: "678908989",
			Language:    "fr",
			PinCode:     otps[len(otps)-1].PinCode,
		})
		if code != http.StatusBadRequest {
			t.Fatalf("CheckOTP() status code = %d; want = %d", code, http.StatusBadRequest)
		}
	})

	t.Run("return 400 if the user otp is not correct", func(t *testing.T) {
		Init()

		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		code, _, _ := helpertest.MakePostRequest(mux, "/otp/check", helpertest.CreateFormHeader(), handlers.CheckOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
			PinCode:     "0000",
		})
		if code != http.StatusBadRequest {
			t.Fatalf("CheckOTP() status code = %d; want = %d", code, http.StatusBadRequest)
		}
	})

	t.Run("return 400 if there is any otp active after checking", func(t *testing.T) {
		Init()

		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		_, _, _ = helpertest.MakePostRequest(mux, "/otp/check", helpertest.CreateFormHeader(), handlers.CheckOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
			PinCode:     otps[len(otps)-1].PinCode,
		})

		var activeOTPExists bool = false
		for i := range otps {
			otp := otps[i]
			if otp.Active {
				activeOTPExists = true
			}
		}
		if activeOTPExists {
			t.Fatalf("CheckOTP() status active-otp-exists = %t; want = %t", activeOTPExists, false)
		}
	})

	t.Run("contains jwt cookies into header", func(t *testing.T) {
		Init()

		_, _, _ = helpertest.MakePostRequest(mux, "/otp", helpertest.CreateFormHeader(), handlers.CreateOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
		})

		_, _, _ = helpertest.MakePostRequest(mux, "/otp/check", helpertest.CreateFormHeader(), handlers.CheckOTPRequest{
			PhoneNumber: "695165033",
			Language:    "fr",
			PinCode:     otps[len(otps)-1].PinCode,
		})
		// headers.
	})
}
