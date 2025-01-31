package storage

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"stockinos.com/api/models"
)

type CreateOrganizationParams struct {
	Name    string
	Bio     string
	Email   string
	Address models.Address

	CreatedBy       primitive.ObjectID
	OwnedBy         primitive.ObjectID
	InvitationToken string
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (*models.Organization, error) {
	var organization models.Organization = models.Organization{
		Id:    primitive.NewObjectID(),
		Name:  arg.Name,
		Bio:   arg.Bio,
		Email: arg.Email,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		CreatedBy:       arg.CreatedBy,
		OwnedBy:         arg.OwnedBy,
		InvitationToken: arg.InvitationToken,
	}

	_, err := q.organizationsCollection.InsertOne(ctx, organization)
	if err != nil {
		return nil, err
	} else {
		return &organization, nil
	}
}

type GetAllCompaniesParams struct {
	UserId primitive.ObjectID
}

func (q *Queries) GetAllCompanies(ctx context.Context, arg GetAllCompaniesParams) ([]*models.Organization, error) {
	var organizations []*models.Organization

	filter := bson.M{
		"created_by": arg.UserId,
		"deleted_at": nil,
	}
	cursor, err := q.organizationsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &organizations); err != nil {
		return nil, err
	}
	return organizations, nil
}

type GetOrganizationParams struct {
	Id primitive.ObjectID
}

func (q *Queries) GetOrganization(ctx context.Context, arg GetOrganizationParams) (*models.Organization, error) {
	var organization models.Organization

	filter := bson.M{
		"_id":        arg.Id,
		"deleted_at": nil,
	}
	err := q.organizationsCollection.FindOne(ctx, filter).Decode(&organization)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &organization, nil
}

type GetOrganizationFromInvitationTokenParams struct {
	InvitationToken string
}

func (q *Queries) GetOrganizationFromInvitationToken(ctx context.Context, arg GetOrganizationFromInvitationTokenParams) (*models.Organization, error) {
	var organization models.Organization

	filter := bson.M{
		"invitation_token": arg.InvitationToken,
		"deleted_at":       nil,
	}
	err := q.organizationsCollection.FindOne(ctx, filter).Decode(&organization)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &organization, nil
}

type UpdateOrganizationParams struct {
	Id   primitive.ObjectID
	Name string
	Bio  string
}

func (q *Queries) UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (*models.Organization, error) {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"name":        arg.Name,
			"description": arg.Bio,
		},
	}
	after := options.After

	var organization models.Organization
	err := q.organizationsCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&organization)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &organization, nil
}

type DeleteOrganizationParams struct {
	Id primitive.ObjectID
}

func (q *Queries) DeleteOrganization(ctx context.Context, arg DeleteOrganizationParams) error {
	filter := bson.M{
		"_id":        arg.Id,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	after := options.After

	var organization models.Organization
	err := q.organizationsCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&organization)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
