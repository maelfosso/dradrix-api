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

type CreateDataParams struct {
	Values     map[string]interface{}
	ActivityId primitive.ObjectID
	CreatedBy  primitive.ObjectID
}

func (q *Queries) CreateData(ctx context.Context, arg CreateDataParams) (*models.Data, error) {
	var data models.Data = models.Data{
		Id:     primitive.NewObjectID(),
		Values: arg.Values,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		// DeletedAt:   null.,

		ActivityId: arg.ActivityId,
		CreatedBy:  arg.CreatedBy,
	}

	_, err := q.datasCollections.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	} else {
		return &data, nil
	}
}

type GetDataParams struct {
	Id primitive.ObjectID
}

func (q *Queries) GetData(ctx context.Context, arg GetDataParams) (*models.Data, error) {
	var data models.Data

	filter := bson.M{
		"_id":        arg.Id,
		"deleted_at": nil,
	}
	err := q.datasCollections.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &data, nil
}

type GetAllDataFromActivityParams struct {
	ActivityId primitive.ObjectID
}

func (q *Queries) GetAllDataFromActivity(ctx context.Context, arg GetAllDataFromActivityParams) ([]*models.Data, error) {
	var data []*models.Data

	filter := bson.M{
		"activity_id": arg.ActivityId,
		"deleted_at":  nil,
	}
	cursor, err := q.datasCollections.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

type DeleteDataParams struct {
	Id primitive.ObjectID
}

func (q *Queries) DeleteData(ctx context.Context, arg DeleteDataParams) error {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	_, err := q.datasCollections.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}
