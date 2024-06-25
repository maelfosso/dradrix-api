package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/fatih/structs"
	"github.com/go-faker/faker/v4"
	"stockinos.com/api/handlers"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
)

func SendPostRequest(url, body string) (*http.Request, *http.Response, error) {
	bodyReader := bytes.NewReader([]byte(body))

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return req, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		// fmt.Printf("client: error making http request: %s\n", err)
		return req, res, err
	}

	return req, res, nil
}

func getResponseData(w *http.Response) string {
	bodyBytes, err := io.ReadAll(w.Body)
	if err != nil {
		panic(err)
	}
	return strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func extractResponseData(response string, data interface{}) {
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		log.Println("extractResponseData : ", response, err.Error())
		panic(err)
	}
}

func TestAuth(t *testing.T) {
	t.Run("one otp request", func(t *testing.T) {

		integrationtest.SkipifShort(t)
		ctx := context.Background()

		stopServer, db := integrationtest.CreateServer()
		defer stopServer()

		host := "http://localhost:8081"
		phoneNumber := faker.Phonenumber()

		// Sign in
		body := fmt.Sprintf(`{
			"phone_number": "%s",
			"language": "en"
		}`, phoneNumber)

		_, w, err := SendPostRequest(fmt.Sprintf("%s/auth/otp", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		var dataResponse handlers.CreateOTPResponse
		extractResponseData(getResponseData(w), &dataResponse)
		if dataResponse.PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: auth/otp: response - got %s; want %s", dataResponse.PhoneNumber, phoneNumber)
		}

		// Check DB
		// 1- Check if the OTP is created; that there is only one activated OTP from that phone number
		otps, err := db.Storage.GetAllOTPs(ctx, storage.GetAllOTPsParams{
			PhoneNumber: phoneNumber,
		})
		if err != nil {
			t.Fatalf("Auth e2e: GetAllOTPs(): error - got %v; want nil", err.Error())
		}
		if len(otps) != 1 {
			t.Fatalf("Auth e2e: GetAllOTPs(): number of otps in db - got %d; want 1", len(otps))
		}
		if otps[0].PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: GetAllOTPs(): otp created has the wrong phone number - got %s; want %s", otps[0].PhoneNumber, phoneNumber)
		}
		if !otps[0].Active {
			t.Fatalf("Auth e2e: GetAllOTPs(): otp created is not activated - got %v; want true", otps[0].Active)
		}

		var activatedOTP *models.OTP
		nActivatedOTPs := 0
		for _, otp := range otps {
			if otp.Active {
				nActivatedOTPs++
				activatedOTP = otp
			}
		}
		if nActivatedOTPs != 1 {
			t.Fatalf("Auth e2e: GetAllOTPs(): number of activated otps in db - got %d; want 1", nActivatedOTPs)
		}

		// 2- Check if the user is created
		users, err := db.Storage.GetAllUsers(ctx, storage.GetAllUsersParams{})
		if err != nil {
			t.Fatalf("Auth e2e: GetAllUsers(): error - got %v; want nil", err.Error())
		}
		if len(users) != 1 {
			t.Fatalf("Auth e2e: GetAllUsers(): number of users in db - got %d; want 1", len(users))
		}
		if users[0].PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: GetAllUsers(): user created has the wrong phone number - got %s; want %s", users[0].PhoneNumber, phoneNumber)
		}

		// Check OTP with wrong phone number
		body = `{
		"phone_number": "0000-0000 223",
		"language": "en",
		"pin_code": "000Z"
	}`

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		errorWrongPinCode := getResponseData(w)
		wantResponse := "ERR_AUTH_CHK_OTP_02"
		if errorWrongPinCode != wantResponse {
			t.Fatalf("Auth e2e: auth/otp/check: response wrong phone number - got %s; want %s", errorWrongPinCode, wantResponse)
		}

		// Check OTP with wrong pin-code
		body = fmt.Sprintf(`{
		"phone_number": "%s",
		"language": "en",
		"pin_code": "000Z"
	}`, phoneNumber)

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		errorWrongPinCode = getResponseData(w)
		wantResponse = "ERR_AUTH_CHK_OTP_03"
		if errorWrongPinCode != wantResponse {
			t.Fatalf("Auth e2e: auth/otp/check: response wrong pin code - got %s; want %s", errorWrongPinCode, wantResponse)
		}

		body = fmt.Sprintf(`{
		"phone_number": "%s",
		"language": "en",
		"pin_code": "%s"
	}`, phoneNumber, activatedOTP.PinCode)

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		var checkOTPResponse handlers.CheckOTPResponse
		extractResponseData(getResponseData(w), &checkOTPResponse)
		if checkOTPResponse.User.PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: auth/otp/check: response - got %s; want %s", checkOTPResponse.User.PhoneNumber, phoneNumber)
		}

		cookies := w.Cookies()
		var foundJwtCookie *http.Cookie
		foundJwtCookie = nil
		for _, c := range cookies {
			if c.Name == "jwt" {
				foundJwtCookie = c
			}
		}
		if foundJwtCookie == nil {
			t.Fatalf("Auth e2e: auth/otp/check: jwt cookie not found - got nil; want not nil")
		}
		jwtToken, err := services.GenerateJwtToken(structs.Map(&checkOTPResponse.User))
		if err != nil {
			t.Fatalf("Auth e2e: Generate JWT Token error: got %v; want nil", err.Error())
		}
		if foundJwtCookie.Value != jwtToken {
			t.Fatalf("Auth e2e: auth/opt/check: wrong jwt cookie: got %s; want %s", foundJwtCookie.Value, jwtToken)
		}
	})

	t.Run("multiple otp requests", func(t *testing.T) {

		integrationtest.SkipifShort(t)
		ctx := context.Background()

		stopServer, db := integrationtest.CreateServer()
		defer stopServer()

		host := "http://localhost:8081"
		phoneNumber := faker.Phonenumber()

		// 1st OTP Request
		body := fmt.Sprintf(`{
			"phone_number": "%s",
			"language": "en"
		}`, phoneNumber)

		_, w, err := SendPostRequest(fmt.Sprintf("%s/auth/otp", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: 1st otp request: error - got %+v; want nil", err.Error())
		}

		var dataResponse handlers.CreateOTPResponse
		extractResponseData(getResponseData(w), &dataResponse)
		if dataResponse.PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: auth/otp: 1st otp request: response - got %s; want %s", dataResponse.PhoneNumber, phoneNumber)
		}

		// 2nd OTP Request
		body = fmt.Sprintf(`{
			"phone_number": "%s",
			"language": "en"
		}`, phoneNumber)

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: 2nd otp request: error - got %+v; want nil", err.Error())
		}

		extractResponseData(getResponseData(w), &dataResponse)
		if dataResponse.PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: auth/otp: 2nd otp request: response - got %s; want %s", dataResponse.PhoneNumber, phoneNumber)
		}

		// 1- Check if the OTP is created; that there is only one activated OTP from that phone number
		otps, err := db.Storage.GetAllOTPs(ctx, storage.GetAllOTPsParams{
			PhoneNumber: phoneNumber,
		})
		if err != nil {
			t.Fatalf("Auth e2e: GetAllOTPs(): error - got %v; want nil", err.Error())
		}
		if len(otps) != 2 {
			t.Fatalf("Auth e2e: GetAllOTPs(): number of otps in db - got %d; want 2", len(otps))
		}
		for ix, otp := range otps {
			if otp.PhoneNumber != phoneNumber {
				t.Fatalf("Auth e2e: GetAllOTPs(): otp(%d) created has the wrong phone number - got %s; want %s", ix, otp.PhoneNumber, phoneNumber)
			}
			if ix < (len(otps)-1) && otp.Active {
				t.Fatalf("Auth e2e: GetAllOTPs(): otp(%d) created is activated - got %v; want false", ix, otp.Active)
			}
			if ix == (len(otps)-1) && !otp.Active {
				t.Fatalf("Auth e2e: GetAllOTPs(): otp(%d) created is not activated - got %v; want true", ix, otp.Active)
			}
		}

		var activatedOTP *models.OTP
		nActivatedOTPs := 0
		for _, otp := range otps {
			if otp.Active {
				nActivatedOTPs++
				activatedOTP = otp
			}
		}
		if nActivatedOTPs != 1 {
			t.Fatalf("Auth e2e: GetAllOTPs(): number of activated otps in db - got %d; want 1", nActivatedOTPs)
		}

		// 2- Check if the user is created
		users, err := db.Storage.GetAllUsers(ctx, storage.GetAllUsersParams{})
		if err != nil {
			t.Fatalf("Auth e2e: GetAllUsers(): error - got %v; want nil", err.Error())
		}
		if len(users) != 1 {
			t.Fatalf("Auth e2e: GetAllUsers(): number of users in db - got %d; want 1", len(users))
		}
		if users[0].PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: GetAllUsers(): user created has the wrong phone number - got %s; want %s", users[0].PhoneNumber, phoneNumber)
		}

		// Check OTP with wrong phone number
		body = `{
			"phone_number": "0000-0000 223",
			"language": "en",
			"pin_code": "000Z"
		}`

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		errorWrongPinCode := getResponseData(w)
		wantResponse := "ERR_AUTH_CHK_OTP_02"
		if errorWrongPinCode != wantResponse {
			t.Fatalf("Auth e2e: auth/otp/check: response wrong phone number - got %s; want %s", errorWrongPinCode, wantResponse)
		}

		// Check OTP with wrong pin-code
		body = fmt.Sprintf(`{
		"phone_number": "%s",
		"language": "en",
		"pin_code": "000Z"
	}`, phoneNumber)

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		errorWrongPinCode = getResponseData(w)
		wantResponse = "ERR_AUTH_CHK_OTP_03"
		if errorWrongPinCode != wantResponse {
			t.Fatalf("Auth e2e: auth/otp/check: response wrong pin code - got %s; want %s", errorWrongPinCode, wantResponse)
		}

		body = fmt.Sprintf(`{
		"phone_number": "%s",
		"language": "en",
		"pin_code": "%s"
	}`, phoneNumber, activatedOTP.PinCode)

		_, w, err = SendPostRequest(fmt.Sprintf("%s/auth/otp/check", host), body)
		if err != nil {
			t.Fatalf("Auth e2e: auth/otp: error - got %+v; want nil", err.Error())
		}

		var checkOTPResponse handlers.CheckOTPResponse
		extractResponseData(getResponseData(w), &checkOTPResponse)
		if checkOTPResponse.User.PhoneNumber != phoneNumber {
			t.Fatalf("Auth e2e: auth/otp/check: response - got %s; want %s", checkOTPResponse.User.PhoneNumber, phoneNumber)
		}

		cookies := w.Cookies()
		var foundJwtCookie *http.Cookie
		foundJwtCookie = nil
		for _, c := range cookies {
			if c.Name == "jwt" {
				foundJwtCookie = c
			}
		}
		if foundJwtCookie == nil {
			t.Fatalf("Auth e2e: auth/otp/check: jwt cookie not found - got nil; want not nil")
		}
		jwtToken, err := services.GenerateJwtToken(structs.Map(&checkOTPResponse.User))
		if err != nil {
			t.Fatalf("Auth e2e: Generate JWT Token error: got %v; want nil", err.Error())
		}
		if foundJwtCookie.Value != jwtToken {
			t.Fatalf("Auth e2e: auth/opt/check: wrong jwt cookie: got %s; want %s", foundJwtCookie.Value, jwtToken)
		}
	})
}
