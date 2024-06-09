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

type CreateCompanyParams struct {
	Name        string
	Description string

	CreatedBy primitive.ObjectID
}

func (q *Queries) CreateCompany(ctx context.Context, arg CreateCompanyParams) (*models.Company, error) {
	var company models.Company = models.Company{
		Id:          primitive.NewObjectID(),
		Name:        arg.Name,
		Description: arg.Description,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		CreatedBy: arg.CreatedBy,
	}

	_, err := q.companiesCollection.InsertOne(ctx, company)
	if err != nil {
		return nil, err
	} else {
		return &company, nil
	}
}

type GetCompanyParams struct {
	Id primitive.ObjectID
}

func (q *Queries) GetCompany(ctx context.Context, arg GetCompanyParams) (*models.Company, error) {
	var company models.Company

	filter := bson.M{
		"_id":        arg.Id,
		"deleted_at": nil,
	}
	err := q.companiesCollection.FindOne(ctx, filter).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &company, nil
}

type UpdateCompanyParams struct {
	Id          primitive.ObjectID
	Name        string
	Description string
}

func (q *Queries) UpdateCompany(ctx context.Context, arg UpdateCompanyParams) (*models.Company, error) {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"name":        arg.Name,
			"description": arg.Description,
		},
	}
	after := options.After

	var company models.Company
	err := q.companiesCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &company, nil
}

type DeleteCompanyParams struct {
	Id primitive.ObjectID
}

func (q *Queries) DeleteCompany(ctx context.Context, arg DeleteCompanyParams) error {
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

	var company models.Company
	err := q.companiesCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&company)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
