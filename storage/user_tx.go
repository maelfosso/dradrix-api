package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func (store *MongoStorage) AddUserInCompanyTx(ctx context.Context) {
	store.withTx(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		store.GetUserByUsername(sessCtx, GetUserByUsernameParams{
			PhoneNumber: "68902323",
		})

		return nil, nil
	})
}
