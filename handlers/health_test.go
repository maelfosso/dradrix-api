package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/handlers"
	"stockinos.com/api/helpertest"
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
		code, _, _ := helpertest.MakeGetRequest(mux, "/health")
		if code != http.StatusOK {
			t.Fatalf("Health() status = %d; want = %d", code, http.StatusOK)
		}
	})

	t.Run("returns 502 if the database connot be pinged", func(t *testing.T) {
		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{err: errors.New("Oh, no!")})
		code, _, _ := helpertest.MakeGetRequest(mux, "/health")
		if code != http.StatusBadGateway {
			t.Fatalf("Health() status = %d; want = %d", code, http.StatusBadGateway)
		}
	})
}
