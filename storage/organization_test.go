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

func TestOrganization(t *testing.T) {
	db, disconnect := integrationtest.CreateDatabase()
	defer disconnect()

	tests := map[string]func(*testing.T, *storage.Database){
		"CreateOrganization": testCreateOrganization,
		"UpdateOrganization": testUpdateOrganization,
		"DeleteOrganization": testDeleteOrganization,
		"GetAllCompanies":    testGetAllCompanies,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			integrationtest.Cleanup(db)

			tc(t, db)
		})
	}
}

func testCreateOrganization(t *testing.T, db *storage.Database) {
	organization, err := db.Storage.CreateOrganization(context.Background(), storage.CreateOrganizationParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})
	if err != nil {
		t.Fatalf("CreateOrganization() error %+v", err)
	}
	if organization == nil {
		t.Fatalf("CreateOrganization() organization is nil; want not nil")
	}

	got, err := db.Storage.GetOrganization(context.Background(), storage.GetOrganizationParams{
		Id: organization.Id,
	})
	if err != nil {
		t.Fatalf("GetOrganization() error %+v", err)
	}
	if err := organizationEq(got, organization); err != nil {
		t.Fatalf("GetOrganization() %v", err)
	}
}

func testGetAllCompanies(t *testing.T, db *storage.Database) {
	const NUM_COMPANIES_CREATED = 3
	userA := primitive.NewObjectID()
	userB := primitive.NewObjectID()

	var organizations []*models.Organization

	for i := 0; i < NUM_COMPANIES_CREATED; i++ {
		organization, _ := db.Storage.CreateOrganization(context.Background(), storage.CreateOrganizationParams{
			Name:        sfaker.Company().Name(),
			Description: gofaker.Paragraph(),

			CreatedBy: userA,
		})
		organizations = append(organizations, organization)
	}

	got, err := db.Storage.GetAllCompanies(context.Background(), storage.GetAllCompaniesParams{
		UserId: userA,
	})
	if err != nil {
		t.Fatalf("GetAllCompanies(): got error; want nil")
	}
	if len(got) != NUM_COMPANIES_CREATED {
		t.Fatalf("GetAllCompanies(): got %d organizations; want %d organizations", len(got), NUM_COMPANIES_CREATED)
	}
	if err := organizationEq(got[0], organizations[0]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}
	if err := organizationEq(got[1], organizations[1]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}
	if err := organizationEq(got[len(got)-1], organizations[len(organizations)-1]); err != nil {
		t.Fatalf("GetAllCompanies(): %v", err)
	}

	got, err = db.Storage.GetAllCompanies(context.Background(), storage.GetAllCompaniesParams{
		UserId: userB,
	})
	if err != nil {
		t.Fatalf("GetAllCompanies(): got error; want nil")
	}
	if len(got) != 0 {
		t.Fatalf("GetAllCompanies(): got %d organizations; want %d organizations", len(got), 0)
	}
}

func testUpdateOrganization(t *testing.T, db *storage.Database) {
	organization, _ := db.Storage.CreateOrganization(context.Background(), storage.CreateOrganizationParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})

	got, err := db.Storage.UpdateOrganization(context.Background(), storage.UpdateOrganizationParams{
		Id:          organization.Id,
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),
	})
	if err != nil {
		t.Fatalf("UpdateOrganization(): got error=%+v; want not error", err)
	}
	if organization.Name == got.Name {
		t.Fatalf("UpdateOrganization(): got same name (%s); want different name", got.Name)
	}
	if organization.Description == got.Description {
		t.Fatalf("UpdateOrganization() Name: got same description (%s); want different description", got.Description)
	}
	if organization.CreatedBy != got.CreatedBy {
		t.Fatalf("UpdateOrganization() Created by: got %s; want %s", got.CreatedBy.String(), organization.CreatedBy.String())
	}

	got, err = db.Storage.GetOrganization(context.Background(), storage.GetOrganizationParams{
		Id: organization.Id,
	})
	if err != nil {
		t.Fatalf("GetOrganization(): got error=%+v; want not error", err)
	}
	if organization.Name == got.Name {
		t.Fatalf("GetOrganization(): got same name (%s); want different name", got.Name)
	}
	if organization.Description == got.Description {
		t.Fatalf("GetOrganization() Name: got same description (%s); want different description", got.Description)
	}
	if organization.CreatedBy != got.CreatedBy {
		t.Fatalf("GetOrganization() Created by: got %s; want %s", got.CreatedBy.String(), organization.CreatedBy.String())
	}
}

func testDeleteOrganization(t *testing.T, db *storage.Database) {
	organization, _ := db.Storage.CreateOrganization(context.Background(), storage.CreateOrganizationParams{
		Name:        sfaker.Company().Name(),
		Description: gofaker.Paragraph(),

		CreatedBy: primitive.NewObjectID(),
	})

	err := db.Storage.DeleteOrganization(context.Background(), storage.DeleteOrganizationParams{
		Id: organization.Id,
	})
	if err != nil {
		t.Fatalf("DeleteOrganization() error %+v", err)
	}

	got, err := db.Storage.GetOrganization(context.Background(), storage.GetOrganizationParams{
		Id: organization.Id,
	})
	if err != nil {
		t.Fatalf("GetOrganization() error %+v", err)
	}
	if got != nil {
		t.Errorf("GetOrganization() got a organization (%s); want nil", got.Name)
	}
}

func organizationEq(got, want *models.Organization) error {
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
