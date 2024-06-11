package storage_test

import (
	"context"
	"fmt"
	"testing"

	gofaker "github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"stockinos.com/api/integrationtest"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
	sfaker "syreclabs.com/go/faker"
)

func TestCompany(t *testing.T) {
	db, disconnect := integrationtest.CreateDatabase()
	defer disconnect()

	tests := map[string]func(*testing.T, *storage.Database){
		"CreateCompany":   testCreateCompany,
		"UpdateCompany":   testUpdateCompany,
		"DeleteCompany":   testDeleteCompany,
		"GetAllCompanies": testGetAllCompanies,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			integrationtest.Cleanup(db)

			tc(t, db)
		})
	}
}

func testCreateCompany(t *testing.T, db *storage.Database) {
	company, err := db.Storage.CreateCompany(context.Background(), storage.CreateCompanyParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})
	if err != nil {
		t.Fatalf("CreateCompany() error %+v", err)
	}
	if company == nil {
		t.Fatalf("CreateCompany() company is nil; want not nil")
	}

	got, err := db.Storage.GetCompany(context.Background(), storage.GetCompanyParams{
		Id: company.Id,
	})
	if err != nil {
		t.Fatalf("GetCompany() error %+v", err)
	}
	if err := companyEq(got, company); err != nil {
		t.Fatalf("GetCompany() %v", err)
	}
}

func testGetAllCompanies(t *testing.T, db *storage.Database) {
	const NUM_COMPANIES_CREATED = 3
	userA := primitive.NewObjectID()
	userB := primitive.NewObjectID()

	var companies []*models.Company

	for i := 0; i < NUM_COMPANIES_CREATED; i++ {
		company, _ := db.Storage.CreateCompany(context.Background(), storage.CreateCompanyParams{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),

			CreatedBy: userA,
		})
		companies = append(companies, company)
	}

	got, err := db.Storage.GetAllCompanies(context.Background(), storage.GetAllCompaniesParams{
		UserId: userA,
	})
	if err != nil {
		t.Fatalf("GetAllCompanies(): got error; want nil")
	}
	if len(got) != NUM_COMPANIES_CREATED {
		t.Fatalf("GetAllCompanies(): got %d companies; want %d companies", len(got), NUM_COMPANIES_CREATED)
	}
	if err := companyEq(got[0], companies[0]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}
	if err := companyEq(got[1], companies[1]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}
	if err := companyEq(got[len(got)-1], companies[len(companies)-1]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}

	got, err = db.Storage.GetAllCompanies(context.Background(), storage.GetAllCompaniesParams{
		UserId: userB,
	})
	if err != nil {
		t.Fatalf("GetAllCompanies(): got error; want nil")
	}
	if len(got) != 0 {
		t.Fatalf("GetAllCompanies(): got %d companies; want %d companies", len(got), 0)
	}
}

func testUpdateCompany(t *testing.T, db *storage.Database) {
	company, _ := db.Storage.CreateCompany(context.Background(), storage.CreateCompanyParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})

	got, err := db.Storage.UpdateCompany(context.Background(), storage.UpdateCompanyParams{
		Id:          company.Id,
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),
	})
	if err != nil {
		t.Fatalf("UpdateCompany(): got error=%+v; want not error", err)
	}
	if company.Name == got.Name {
		t.Fatalf("UpdateCompany(): got same name (%s); want different name", got.Name)
	}
	if company.Description == got.Description {
		t.Fatalf("UpdateCompany() Name: got same description (%s); want different description", got.Description)
	}
	if company.CreatedBy != got.CreatedBy {
		t.Fatalf("UpdateCompany() Created by: got %s; want %s", got.CreatedBy.String(), company.CreatedBy.String())
	}

	got, err = db.Storage.GetCompany(context.Background(), storage.GetCompanyParams{
		Id: company.Id,
	})
	if err != nil {
		t.Fatalf("GetCompany(): got error=%+v; want not error", err)
	}
	if company.Name == got.Name {
		t.Fatalf("GetCompany(): got same name (%s); want different name", got.Name)
	}
	if company.Description == got.Description {
		t.Fatalf("GetCompany() Name: got same description (%s); want different description", got.Description)
	}
	if company.CreatedBy != got.CreatedBy {
		t.Fatalf("GetCompany() Created by: got %s; want %s", got.CreatedBy.String(), company.CreatedBy.String())
	}
}

func testDeleteCompany(t *testing.T, db *storage.Database) {
	company, _ := db.Storage.CreateCompany(context.Background(), storage.CreateCompanyParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})

	err := db.Storage.DeleteCompany(context.Background(), storage.DeleteCompanyParams{
		Id: company.Id,
	})
	if err != nil {
		t.Fatalf("DeleteCompany() error %+v", err)
	}

	got, err := db.Storage.GetCompany(context.Background(), storage.GetCompanyParams{
		Id: company.Id,
	})
	if err != nil {
		t.Fatalf("GetCompany() error %+v", err)
	}
	if got != nil {
		t.Errorf("GetCompany() got a company (%s); want nil", got.Name)
	}
}

func companyEq(got, want *models.Company) error {
	if got == want {
		return nil
	}
	if got == nil {
		return fmt.Errorf("got nil; want %v", want)
	}
	if want == nil {
		return fmt.Errorf("got %v; want nil", got)
	}
	if got.Id != want.Id {
		return fmt.Errorf("got.Id = %s; want %s", got.Id, want.Id)
	}
	if got.Name != want.Name {
		return fmt.Errorf("got.Name = %s; want %s", got.Name, want.Name)
	}
	if got.Description != want.Description {
		return fmt.Errorf("got.Description = %s; want %s", got.Description, want.Description)
	}
	if got.CreatedBy != want.CreatedBy {
		return fmt.Errorf("got.CreatedBy = %s; want %s", got.CreatedBy, want.CreatedBy)
	}

	return nil
}
