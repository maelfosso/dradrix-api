package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type CreateDataParams struct {
	Values     map[string]any
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
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID
}

func (q *Queries) GetData(ctx context.Context, arg GetDataParams) (*models.Data, error) {
	var data models.Data

	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
		"deleted_at":  nil,
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

type GetAllDataParams struct {
	ActivityId primitive.ObjectID
}

func (q *Queries) GetAllData(ctx context.Context, arg GetAllDataParams) ([]*models.Data, error) {
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
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID
}

func (q *Queries) DeleteData(ctx context.Context, arg DeleteDataParams) error {
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
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

type UpdateSetInDataParams struct {
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID

	Field string
	Value interface{}
}

func (q *Queries) UpdateSetInData(ctx context.Context, arg UpdateSetInDataParams) (*models.Data, error) {
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$set": bson.M{
			"values": bson.M{
				arg.Field: arg.Value,
			},
		},
	}

	return CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
}

type UpdateAddToDataParams struct {
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID

	Position uint
	Field    string
	Value    interface{}
}

func (q *Queries) UpdateAddToData(ctx context.Context, arg UpdateAddToDataParams) (*models.Data, error) {
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$push": bson.M{
			"values": bson.M{
				arg.Field: bson.M{
					"$each": bson.A{
						arg.Value,
					},
					"$position": arg.Position,
				},
			},
		},
	}

	return CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
}

type UpdateRemoveFromDataParams struct {
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID

	Position uint
	Field    string
	// Value interface{}
}

func (q *Queries) UpdateRemoveFromData(ctx context.Context, arg UpdateRemoveFromDataParams) (*models.Data, error) {
	field := fmt.Sprintf("%s.%d", arg.Field, arg.Position)
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$pop": bson.M{
			"values": bson.M{
				field: 1,
			},
		},
	}

	return CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
}
