package helpertest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"unicode"
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

	return result.StatusCode, result.Header, strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func MakePostRequest(handler http.Handler, target string, header http.Header, data interface{}) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, target, bytes.NewReader(body))

	req.Header = header
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func MakePutRequest(handler http.Handler, target string, header http.Header, data interface{}) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPut, target, bytes.NewReader(body))

	req.Header = header
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func MakeDeleteRequest(handler http.Handler, target string, header http.Header, data interface{}) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodDelete, target, bytes.NewReader(body))

	req.Header = header
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func CreateFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return header
}
