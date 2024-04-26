package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/handlers"
)

type pingerMock struct {
	err error
}

func (p *pingerMock) Ping(ctx context.Context) error {
	return p.err
}

func TestHealth(t *testing.T) {
	t.Run("returns 200", func(t *testing.T) {
		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{})
		code, _, _ := makeGetRequest(mux, "/health")
		if code != http.StatusOK {
			t.Fatalf("Health() status = %d; want = %d", code, http.StatusOK)
		}
	})

	t.Run("returns 502 if the database connot be pinged", func(t *testing.T) {
		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{err: errors.New("Oh, no!")})
		code, _, _ := makeGetRequest(mux, "/health")
		if code != http.StatusBadGateway {
			t.Fatalf("Health() status = %d; want = %d", code, http.StatusBadGateway)
		}
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
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
