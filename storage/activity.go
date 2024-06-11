package storage

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type CreateActivityParams struct {
	Name        string
	Description string
	Fields      []models.ActivityFields

	CompanyId primitive.ObjectID
	CreatedBy primitive.ObjectID
}

func (q *Queries) CreateActivity(ctx context.Context, arg CreateActivityParams) (*models.Activity, error) {
	var activity models.Activity = models.Activity{
		Id:          primitive.NewObjectID(),
		Name:        arg.Name,
		Description: arg.Description,
		Fields:      arg.Fields,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		// DeletedAt:   null.,

		CompanyId: arg.CompanyId,
		CreatedBy: arg.CreatedBy,
	}

	_, err := q.activitiesCollection.InsertOne(ctx, activity)
	if err != nil {
		return nil, err
	} else {
		return &activity, nil
	}
}

type GetActivityFromCompanyParams struct {
	Id        primitive.ObjectID
	CompanyId primitive.ObjectID
}

func (q *Queries) GetActivityFromCompany(ctx context.Context, arg GetActivityFromCompanyParams) (*models.Activity, error) {
	var activity models.Activity

	filter := bson.M{
		"_id":        arg.Id,
		"company_id": arg.CompanyId,
		"deleted_at": nil,
	}
	err := q.activitiesCollection.FindOne(ctx, filter).Decode(&activity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &activity, nil
}

type GetAllActivitiesFromCompanyParams struct {
	CompanyId primitive.ObjectID
}

func (q *Queries) GetAllActivitiesFromCompany(ctx context.Context, arg GetAllActivitiesFromCompanyParams) ([]*models.Activity, error) {
	var activities []*models.Activity

	filter := bson.M{
		"company_id": arg.CompanyId,
		"deleted_at": nil,
	}
	cursor, err := q.activitiesCollection.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

type DeleteActivityParams struct {
	Id        primitive.ObjectID
	CompanyId primitive.ObjectID
}

func (q *Queries) DeleteActivityFromCompany(ctx context.Context, arg DeleteActivityParams) error {
	filter := bson.M{
		"_id":        arg.Id,
		"company_id": arg.CompanyId,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	var activity models.Activity
	err := q.activitiesCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
	).Decode(&activity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
