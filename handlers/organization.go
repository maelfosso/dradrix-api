package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
)

type getOrganizationCtxInterface interface {
	GetOrganization(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error)
}

func (handler *AppHandler) OrganizationMiddleware(mux chi.Router, db getOrganizationCtxInterface) {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			organizationIdParam := chi.URLParamFromCtx(ctx, "organizationId")
			organizationId, err := primitive.ObjectIDFromHex(organizationIdParam)
			if err != nil {
				http.Error(w, "ERR_CMP_MDW_01", http.StatusBadRequest)
				return
			}

			organization, err := db.GetOrganization(ctx, storage.GetOrganizationParams{
				Id: organizationId,
			})
			if err != nil {
				http.Error(w, "ERR_CMP_MDW_02", http.StatusBadRequest)
				return
			}

			if organization == nil {
				http.Error(w, "ERR_CMP_MDW_03", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "organization", organization)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

type getAllCompaniesInterface interface {
	GetAllCompanies(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Organization, error)
}

type GetAllCompaniesResponse struct {
	Companies []*models.Organization `json:"organizations,omitempty"`
}

func (handler *AppHandler) GetAllCompanies(mux chi.Router, db getAllCompaniesInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var organizations []*models.Organization

		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		organizations, err := db.GetAllCompanies(ctx, storage.GetAllCompaniesParams{
			UserId: authUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_GALL_CMP_01", http.StatusBadRequest)
			return
		}

		response := GetAllCompaniesResponse{
			Companies: organizations,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_GALL_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type createOrganizationInterface interface {
	CreateOrganization(ctx context.Context, arg storage.CreateOrganizationParams) (*models.Organization, error)
}

type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateOrganizationResponse struct {
	Organization models.Organization `json:"organization,omitempty"`
}

func (handler *AppHandler) CreateOrganization(mux chi.Router, db createOrganizationInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		var input CreateOrganizationRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		organization, err := db.CreateOrganization(ctx, storage.CreateOrganizationParams{
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   authUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_C_CMP_01", http.StatusBadRequest)
			return
		}

		response := CreateOrganizationResponse{
			Organization: *organization,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_C_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type getOrganizationInterface interface {
	// GetOrganization(ctx context.Context, arg storage.GetOrganizationParams) (*models.Organization, error)
}

type GetOrganizationResponse struct {
	Organization models.Organization `json:"organization,omitempty"`
}

func (handler *AppHandler) GetOrganization(mux chi.Router, db getOrganizationInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organization := ctx.Value("organization").(*models.Organization)

		response := GetOrganizationResponse{
			Organization: *organization,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_C_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type updateOrganizationInterface interface {
	UpdateOrganization(ctx context.Context, arg storage.UpdateOrganizationParams) (*models.Organization, error)
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateOrganizationResponse struct {
	Organization models.Organization `json:"organization,omitempty"`
}

func (handler *AppHandler) UpdateOrganization(mux chi.Router, db updateOrganizationInterface) {
	mux.Put("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input UpdateOrganizationRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		organization := ctx.Value("organization").(*models.Organization)

		updatedOrganization, err := db.UpdateOrganization(ctx, storage.UpdateOrganizationParams{
			Id:          organization.Id,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			http.Error(w, "ERR_U_CMP_01", http.StatusBadRequest)
			return
		}

		response := UpdateOrganizationResponse{
			Organization: *updatedOrganization,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_U_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteOrganizationInterface interface {
	DeleteOrganization(ctx context.Context, arg storage.DeleteOrganizationParams) error
}

type DeleteOrganizationRequest struct{}

type DeleteOrganizationResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

func (handler *AppHandler) DeleteOrganization(mux chi.Router, db deleteOrganizationInterface) {
	mux.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organization := ctx.Value("organization").(*models.Organization)

		err := db.DeleteOrganization(ctx, storage.DeleteOrganizationParams{
			Id: organization.Id,
		})
		if err != nil {
			http.Error(w, "ERR_D_CMP_02", http.StatusBadRequest)
			return
		}

		response := DeleteOrganizationResponse{
			Deleted: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_D_CMP_END", http.StatusBadRequest)
			return
		}
	})
}
