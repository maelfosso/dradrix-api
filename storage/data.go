package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"stockinos.com/api/models"
)

type CreateDataParams struct {
	Values     map[string]any
	ActivityId primitive.ObjectID
	CreatedBy  models.DataAuthor
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

type GetDataFilterByValuesParams struct {
	Values     map[string]any
	ActivityId primitive.ObjectID
}

func (q *Queries) GetDataFilterByValues(ctx context.Context, arg GetDataFilterByValuesParams) (*models.Data, error) {
	var data models.Data

	filter := bson.M{}
	for field, value := range arg.Values {
		filter[fmt.Sprintf("values.%s", field)] = value
	}
	filter["activity_id"] = arg.ActivityId
	filter["deleted_at"] = nil

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
	ActivityId  primitive.ObjectID
	Projections map[string]int
	FilterBy    map[string]any
}

func (q *Queries) GetAllData(ctx context.Context, arg GetAllDataParams) ([]*models.Data, error) {
	var data []*models.Data

	projections := bson.M{}
	if len(arg.Projections) > 0 {
		for key, value := range arg.Projections {
			projections[key] = value
		}
	}
	opts := options.Find().SetProjection(projections)

	filter := bson.M{
		"activity_id": arg.ActivityId,
		"deleted_at":  nil,
	}
	if len(arg.FilterBy) > 0 {
		for key, value := range arg.FilterBy {
			filter[key] = value
		}
	}
	cursor, err := q.datasCollections.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	if data == nil {
		return []*models.Data{}, nil
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
	field := fmt.Sprintf("values.%s", arg.Field)
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$set": bson.M{
			field: arg.Value,
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
	field := fmt.Sprintf("values.%s", arg.Field)
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$push": bson.M{
			field: bson.M{
				"$each": bson.A{
					arg.Value,
				},
				"$position": arg.Position,
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
}

func (q *Queries) UpdateRemoveFromData(ctx context.Context, arg UpdateRemoveFromDataParams) (*models.Data, error) {
	field := fmt.Sprintf("values.%s", arg.Field)
	fieldWithPosition := fmt.Sprintf("%s.%d", field, arg.Position)
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}
	update := bson.M{
		"$unset": bson.M{
			fieldWithPosition: 1,
		},
	}
	_, err := CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
	if err != nil {
		return nil, err
	}

	update = bson.M{
		"$pull": bson.M{
			field: nil,
		},
	}
	return CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
}

type UpdateDataParams struct {
	Id         primitive.ObjectID
	ActivityId primitive.ObjectID

	Values map[string]any
}

func (q *Queries) UpdateData(ctx context.Context, arg UpdateDataParams) (*models.Data, error) {
	filter := bson.M{
		"_id":         arg.Id,
		"activity_id": arg.ActivityId,
	}

	set := bson.M{}
	for key, value := range arg.Values {
		field := fmt.Sprintf("values.%s", key)
		set[field] = value
	}
	update := bson.M{
		"$set": set,
	}

	return CommonUpdateQuery[models.Data](ctx, *q.datasCollections, filter, update)
}
