package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"stockinos.com/api/models"
	"stockinos.com/api/services"
)

type AppHandler struct {
	GetAuthenticatedUser func(r *http.Request) *models.User
	ParsingRequestBody   func(w http.ResponseWriter, r *http.Request, inputs interface{}) (int, error)
}

func NewAppHandler() *AppHandler {
	return &AppHandler{
		GetAuthenticatedUser: func(r *http.Request) *models.User {
			user := r.Context().Value(services.JwtUserKey)
			if user == nil {
				return nil
			}

			if u, ok := user.(*models.User); ok {
				return u
			}

			return nil
		},
		ParsingRequestBody: func(w http.ResponseWriter, r *http.Request, input interface{}) (int, error) {
			// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

			// If the Content-Type header is present, check that it has the value application/json
			contentType := r.Header.Get("Content-Type")
			if contentType != "" {
				mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
				if mediaType != "application/json" {
					// msg := "Content-Type header is not application/json"
					// http.Error(w, "ERR_HDL_PRB_01", http.StatusUnsupportedMediaType)
					return http.StatusUnsupportedMediaType, errors.New("ERR_HDL_PRB_01")
				}
			}

			// Use http.MaxBytesReader to enforce a maximum read of 1MB from the response body.
			// A request body larger that that will now result in Decode() returning a
			// "http: request body too large" error
			r.Body = http.MaxBytesReader(w, r.Body, 1048576)

			// Setup the decodeer and call the DisallowUnknownFields() method on it
			// This will cause Decode to return a "json: unknown field ..." error
			// if it encounters any extra unexpected fields in the JSON.
			// Strictly speaking, it returns an error for "keys which do not match any non-ignored,
			// exported fields in the destination"
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			err := decoder.Decode(&input)
			if err != nil {
				var syntaxError *json.SyntaxError
				var unmarshalTypeError *json.UnmarshalTypeError

				switch {
				// Catch any syntax errors in the JSON and log an error message with
				// the location of the problem
				case errors.As(err, &syntaxError):
					// msg := fmt.Sprintf(
					// 	"Request body contains badly-formed JSON (at position &d)",
					// 	syntaxError.Offset,
					// )
					return http.StatusBadRequest, errors.New("ERR_HDL_PRB_02")

				// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
				// for syntax errors in the JSON
				case errors.Is(err, io.ErrUnexpectedEOF):
					// msg := fmt.Sprintf("Request body contains badly-formed JSON")
					return http.StatusBadRequest, errors.New("ERR_HDL_PRB_03")

				// Catch any type errors, like trying to assign a string in the JSON request body
				// to a int field in our data struct
				case errors.As(err, &unmarshalTypeError):
					// msg := fmt.Sprintf(
					// 	"Request body contains an invalid value for the %q field (at position)",
					// 	unmarshalTypeError.Offset,
					// )
					return http.StatusBadRequest, errors.New("ERR_HDL_PRB_04")

				// Catch error caused by extra unexpected fields in the request body.
				// We extract the field name from the errror message and interpolate it in our custom error message
				case strings.HasPrefix(err.Error(), "json: unknown field "):
					// fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
					// msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
					return http.StatusBadRequest, errors.New("ERR_HDL_PRB_05")

				// An io.EOF error is returned by Decode() if the request body is empty
				case errors.Is(err, io.EOF):
					// msg := "Request body must not be empty"
					return http.StatusBadRequest, errors.New("ERR_HDL_PRB_06")

				// Catch the error caused by the request body being too large.
				case err.Error() == "http: request body too large":
					// msg := "Request body must not be larger than 1MB"
					return http.StatusRequestEntityTooLarge, errors.New("ERR_HDL_PRB_07")

				// Otherwise default to logging the error and sending a 500 internal Server Error response.
				default:
					// msg := err.Error()
					// http.StatusText(http.StatusInternalServerError)
					return http.StatusInternalServerError, errors.New("ERR_HDL_PRB_08")
				}
			}

			// Call decode again, using a pointer to an empty anonymous struct as the destination
			// If the request body contained a single JSON object this will return an io.EOF error.
			// So if we get anything else, we know that there is additional data in the request body.
			err = decoder.Decode(&struct{}{})
			if !errors.Is(err, io.EOF) {
				// msg := "Request body must only contain a single JSON object"
				return http.StatusBadRequest, errors.New("ERR_HDL_PRB_08")
			}

			return -1, nil
		},
	}
}
