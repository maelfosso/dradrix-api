package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type getTeamInterface interface {
	GetMembersFromOrganization(ctx context.Context, arg storage.GetMembersFromOrganizationParams) ([]models.Member, error)
}

type GetTeamResponse struct {
	Members []models.Member
}

func (appHandler *AppHandler) GetTeam(mux chi.Router, db getTeamInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organization := ctx.Value("organization").(*models.Organization)

		members, err := db.GetMembersFromOrganization(ctx, storage.GetMembersFromOrganizationParams{
			OrganizationId: organization.Id,
		})
		if err != nil {
			http.Error(w, "ERR_TEAM_GTM_DB", http.StatusBadRequest)
			return
		}

		response := GetTeamResponse{
			Members: members,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_TEAM_GTM_END", http.StatusBadRequest)
			return
		}
	})
}

func (appHandler *AppHandler) AddMemberTeam(mux chi.Router) {
	mux.Post("", func(w http.ResponseWriter, r *http.Request) {

	})
}

// func (appHandler *AppHandler)
