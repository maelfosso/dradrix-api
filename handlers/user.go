package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
)

type GetCurrentUserResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phone_number"`

	Preferences models.UserPreferences `json:"preferences"`
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
			FirstName:   currentUser.FirstName,
			LastName:    currentUser.LastName,
			Email:       currentUser.Email,
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
