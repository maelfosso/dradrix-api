package helpertest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

func MakeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return result.StatusCode, result.Header, string(bodyBytes)
}

func MakePostRequest(handler http.Handler, target string, header http.Header, body interface{}) (int, http.Header, string) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest(http.MethodPost, target, &buf)

	req.Header = header
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}

func CreateFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return header
}
