package integrationtest

import (
	"log"
	"net/http"
	"testing"
	"time"

	"stockinos.com/api/server"
)

// CreateServer for testing on port 8080, returning a cleanup function stops the server.
// Usage:
//
//	cleanup := CreateServer()
//	defer cleanup()
func CreateServer() func() {
	db, cleanupDB := CreateDatabase()
	s := server.New(server.Options{
		Host:     "localhost",
		Port:     8081,
		Database: db,
	})

	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	for {
		_, err := http.Get("http://localhost:8081/")
		if err == nil {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	return func() {
		if err := s.Stop(); err != nil {
			log.Println("[IntegrationTest] Stop Failed", err)
			panic(err)
		}
		cleanupDB()
	}
}

// SkipIfShort skips t if the "-short" flag is passed to "go test"
func SkipifShort(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}
