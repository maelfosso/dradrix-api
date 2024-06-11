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

type GetActivityParams struct {
	Id        primitive.ObjectID
	CompanyId primitive.ObjectID
}

func (q *Queries) GetActivity(ctx context.Context, arg GetActivityParams) (*models.Activity, error) {
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

type GetAllActivitiesParams struct {
	CompanyId primitive.ObjectID
}

func (q *Queries) GetAllActivities(ctx context.Context, arg GetAllActivitiesParams) ([]*models.Activity, error) {
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

func (q *Queries) DeleteActivity(ctx context.Context, arg DeleteActivityParams) error {
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

type UpdateActivityParams struct {
	Id        primitive.ObjectID
	CompanyId primitive.ObjectID

	Field string
	Value interface{}
	// Type  string
}

func (q *Queries) UpdateActivity(ctx context.Context, arg UpdateActivityParams) (*models.Activity, error) {
	filter := bson.M{
		"_id":        arg.Id,
		"company_id": arg.CompanyId,
	}
	update := bson.M{
		"$set": bson.M{
			arg.Field: arg.Value,
		},
	}
	after := options.After

	var activity models.Activity
	err := q.activitiesCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&activity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &activity, nil
}
