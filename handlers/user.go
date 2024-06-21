package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
)

type GetCurrentUserResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty"`

	Preferences models.UserPreferences `json:"preferences,omitempty"`
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

		response := GetCurrentUserResponse{
			Name:        fmt.Sprintf("%s %s", currentUser.LastName, currentUser.FirstName),
			PhoneNumber: currentUser.PhoneNumber,
			ID:          currentUser.Id,

			Preferences: currentUser.Preferences,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "error encoding the result", http.StatusBadRequest)
			return
		}
	})
}
