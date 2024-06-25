package server_test

import (
	"net/http"
	"testing"

	"stockinos.com/api/integrationtest"
)

func TestServer_Start(t *testing.T) {
	integrationtest.SkipifShort(t)

	t.Run("starts the server and listens for requests", func(t *testing.T) {

		cleanup, _ := integrationtest.CreateServer()
		defer cleanup()

		resp, err := http.Get("http://localhost:8081/")
		if err != nil {
			t.Errorf("Server has not started due to error (%v)", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("CreateServer() status code = %d; want %d", resp.StatusCode, http.StatusOK)
		}
	})
}
