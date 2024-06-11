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

type getAllCompaniesInterface interface {
	GetAllCompanies(ctx context.Context, arg storage.GetAllCompaniesParams) ([]*models.Company, error)
}

type GetAllCompaniesResponse struct {
	Companies []*models.Company `json:"companies,omitempty"`
}

func (handler *AppHandler) GetAllCompanies(mux chi.Router, db getAllCompaniesInterface) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var companies []*models.Company

		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		companies, err := db.GetAllCompanies(ctx, storage.GetAllCompaniesParams{
			UserId: authUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_GALL_CMP_01", http.StatusBadRequest)
			return
		}

		response := GetAllCompaniesResponse{
			Companies: companies,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_GALL_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type getCompanyInterface interface {
	GetCompany(ctx context.Context, arg storage.GetCompanyParams) (*models.Company, error)
}

type GetCompanyResponse struct {
	Company models.Company `json:"company,omitempty"`
}

func (handler *AppHandler) GetCompany(mux chi.Router, db getCompanyInterface) {
	mux.Get("/{companyId}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		companyIdParam := chi.URLParamFromCtx(ctx, "companyId")
		companyId, err := primitive.ObjectIDFromHex(companyIdParam)
		if err != nil {
			http.Error(w, "ERR_GONE_CMP_01", http.StatusBadRequest)
			return
		}

		company, err := db.GetCompany(ctx, storage.GetCompanyParams{
			Id: companyId,
		})
		if err != nil {
			http.Error(w, "ERR_GONE_CMP_02", http.StatusBadRequest)
			return
		}

		response := GetCompanyResponse{
			Company: *company,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_C_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type createCompanyInterface interface {
	CreateCompany(ctx context.Context, arg storage.CreateCompanyParams) (*models.Company, error)
}

type CreateCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCompanyResponse struct {
	Company models.Company `json:"company,omitempty"`
}

func (handler *AppHandler) CreateCompany(mux chi.Router, db createCompanyInterface) {
	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authUser := handler.GetAuthenticatedUser(r)

		var input CreateCompanyRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		company, err := db.CreateCompany(ctx, storage.CreateCompanyParams{
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   authUser.Id,
		})
		if err != nil {
			http.Error(w, "ERR_C_CMP_01", http.StatusBadRequest)
			return
		}

		response := CreateCompanyResponse{
			Company: *company,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_C_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type updateCompanyInterface interface {
	UpdateCompany(ctx context.Context, arg storage.UpdateCompanyParams) (*models.Company, error)
}

type UpdateCompanyRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type UpdateCompanyResponse struct {
	Company models.Company `json:"company,omitempty"`
}

func (handler *AppHandler) UpdateCompany(mux chi.Router, db updateCompanyInterface) {
	mux.Put("/{companyId}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		companyIdParam := chi.URLParamFromCtx(ctx, "companyId")
		companyId, err := primitive.ObjectIDFromHex(companyIdParam)
		if err != nil {
			http.Error(w, "ERR_U_CMP_01", http.StatusBadRequest)
			return
		}

		var input UpdateCompanyRequest
		httpStatus, err := handler.ParsingRequestBody(w, r, &input)
		if err != nil {
			http.Error(w, err.Error(), httpStatus)
			return
		}

		company, err := db.UpdateCompany(ctx, storage.UpdateCompanyParams{
			Id:          companyId,
			Name:        input.Name,
			Description: input.Description,
		})
		if err != nil {
			http.Error(w, "ERR_U_CMP_02", http.StatusBadRequest)
			return
		}

		response := UpdateCompanyResponse{
			Company: *company,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "ERR_U_CMP_END", http.StatusBadRequest)
			return
		}
	})
}

type deleteCompanyInterface interface {
	DeleteCompany(ctx context.Context, arg storage.DeleteCompanyParams) error
}

type DeleteCompanyRequest struct{}

type DeleteCompanyResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

func (handler *AppHandler) DeleteCompany(mux chi.Router, db deleteCompanyInterface) {
	mux.Delete("/{companyId}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		companyIdParam := chi.URLParamFromCtx(ctx, "companyId")
		companyId, err := primitive.ObjectIDFromHex(companyIdParam)
		if err != nil {
			http.Error(w, "ERR_D_CMP_01", http.StatusBadRequest)
			return
		}

		err = db.DeleteCompany(ctx, storage.DeleteCompanyParams{
			Id: companyId,
		})
		if err != nil {
			http.Error(w, "ERR_D_CMP_02", http.StatusBadRequest)
			return
		}

		response := DeleteCompanyResponse{
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
