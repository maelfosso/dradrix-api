package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllMonitoredActivities(mux chi.Router) {
	mux.Get("/monitoring", func(w http.ResponseWriter, r *http.Request) {
		var activities []string

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(activities); err != nil {
			http.Error(w, "ERR_COTP_106", http.StatusBadRequest)
			return
		}
	})
}
