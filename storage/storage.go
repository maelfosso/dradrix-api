package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storage interface {
	Querier
	QuerierTx
}

type MongoStorage struct {
	db *mongo.Client
	*Queries
}

func NewStorage(d Database) Storage {
	return &MongoStorage{
		db:      d.DB,
		Queries: NewQueries(d),
	}
}
