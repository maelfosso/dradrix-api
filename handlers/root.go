// Package handlers contains HTTP handlers used by the server package.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Root(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode("Simple Whatsapp Webhook tester</br>There is no front-end"); err != nil {
			fmt.Println("ERror when encoding - ", err)
			http.Error(w, "error when encoding", http.StatusBadRequest)
			return
		}
	})
}
