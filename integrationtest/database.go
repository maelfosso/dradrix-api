package integrationtest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"stockinos.com/api/storage"
)

var once sync.Once

// CreateDatabase for testing
// Usage:
//
//	db, cleanup := CreateDatabase()
//	defer cleanup()
//	...
func CreateDatabase() (*storage.Database, func()) {

	once.Do(initDatabase)

	db, cleanup := connect("stockinos-test")
	defer cleanup()

	db.DB.Database("stockinos-test").Drop(context.Background())

	return connect("stockinos-test")
}

func initDatabase() {
	db, cleanup := connect("template1")
	defer cleanup()

	for err := db.DB.Ping(context.Background(), nil); err != nil; {
		time.Sleep(100 * time.Millisecond)
	}
}

func connect(name string) (*storage.Database, func()) {
	db := storage.NewDatabase(storage.NewDatabaseOptions{
		URI:  fmt.Sprintf("mongodb://localhost:27017/%s", name),
		Name: name,
	})
	if err := db.Connect(); err != nil {
		panic(err)
	}

	return db, func() {
		if err := db.DB.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func Cleanup(db *storage.Database) {
	db.GetCollection("companies").DeleteMany(context.Background(), bson.M{})
	db.GetCollection("activities").DeleteMany(context.Background(), bson.M{})
	db.GetCollection("otps").DeleteMany(context.Background(), bson.M{})
	db.GetCollection("users").DeleteMany(context.Background(), bson.M{})
}
