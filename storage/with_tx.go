package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func (store *MongoStorage) withTx(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := store.db.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, fn)
	if err != nil {
		return nil, err
	}

	return result, nil
}
