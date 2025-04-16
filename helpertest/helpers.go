package helpertest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"unicode"
)

type ContextData struct {
	Name  string
	Value interface{}
}

func MakeGetRequest(handler http.Handler, target string, ctxData []ContextData) (*http.Request, *http.Response, string) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	for i := 0; i < len(ctxData); i++ {
		cd := ctxData[i]
		if cd.Value != nil {
			req = req.WithContext(context.WithValue(req.Context(), cd.Name, cd.Value))
		}
	}
	handler.ServeHTTP(w, req)

	result := w.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return req, result, strings.TrimFunc(string(bodyBytes), unicode.IsSpace)
}

func MakePostRequest(handler http.Handler, target string, header http.Header, data interface{}, ctxData []ContextData) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	for i := 0; i < len(ctxData); i++ {
		cd := ctxData[i]
		if cd.Value != nil {
			req = req.WithContext(context.WithValue(req.Context(), cd.Name, cd.Value))
		}
	}
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

func MakePutRequest(handler http.Handler, target string, header http.Header, data interface{}, ctxData []ContextData) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPut, target, bytes.NewReader(body))
	for i := 0; i < len(ctxData); i++ {
		cd := ctxData[i]
		if cd.Value != nil {
			req = req.WithContext(context.WithValue(req.Context(), cd.Name, cd.Value))
		}
	}

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

func MakePatchRequest(handler http.Handler, target string, header http.Header, data interface{}, ctxData []ContextData) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPatch, target, bytes.NewReader(body))
	for i := 0; i < len(ctxData); i++ {
		cd := ctxData[i]
		if cd.Value != nil {
			req = req.WithContext(context.WithValue(req.Context(), cd.Name, cd.Value))
		}
	}

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

func MakeDeleteRequest(handler http.Handler, target string, header http.Header, data interface{}, ctxData []ContextData) (int, http.Header, string) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodDelete, target, bytes.NewReader(body))
	for i := 0; i < len(ctxData); i++ {
		cd := ctxData[i]
		if cd.Value != nil {
			req = req.WithContext(context.WithValue(req.Context(), cd.Name, cd.Value))
		}
	}

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
