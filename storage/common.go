package storage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CommonUpdateQuery[T any](ctx context.Context, collection mongo.Collection, filter, update bson.M) (*T, error) {
	after := options.After

	var data T
	err := collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &data, nil
}
