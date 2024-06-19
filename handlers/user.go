package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
)

type GetCurrentUserResult struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

func GetCurrentUser(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		currentUser := ctx.Value(services.JwtUserKey).(*models.User)
		if currentUser == nil {
			log.Println("Error No current user: ")
			http.Error(w, "ERR_NO_CURRENT_USER", http.StatusBadRequest)
			return
		}

		var currentUserResult GetCurrentUserResult
		currentUserResult.Name = currentUser.Name
		currentUserResult.PhoneNumber = currentUser.PhoneNumber
		currentUserResult.ID = currentUser.Id.String()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(currentUserResult); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
