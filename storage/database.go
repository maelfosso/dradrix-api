package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	Storage Storage
	DB      *mongo.Client
	uri     string
	name    string
	log     *zap.Logger
}

type NewDatabaseOptions struct {
	URI  string
	Name string
	Log  *zap.Logger
}

func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	return &Database{
		uri:  opts.URI,
		name: opts.Name,
		log:  opts.Log,
	}
}

func (d *Database) createDataSourceName(withPassword bool) string {
	// password := d.password
	// if !withPassword {
	// 	password = "xxx"
	// }

	// return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
	// 	d.user, password, d.host, d.port, d.name)
	return d.uri
}

func (d *Database) Connect() error {
	d.log.Info("Connecting to database", zap.String("url", d.createDataSourceName(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	// d.DB, err = gorm.Open(postgres.Open(d.createDataSourceName(true)), &gorm.Config{})
	d.DB, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(d.uri),
	)
	if err != nil {
		d.log.Fatal("Failed to connect to the database : ", zap.Error(err))
		return err
	}

	d.Storage = NewStorage(*d)

	err = d.DB.Ping(context.Background(), nil)
	if err != nil {
		d.log.Fatal("Ping to database has failed")
	}

	d.log.Info("Successfully connected to MongoDB")
	return nil
}

func (d *Database) Disconnect() {
	if err := d.DB.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func (d *Database) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := d.DB.Ping(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetCollection(coll string) *mongo.Collection {
	return d.DB.Database(d.name).Collection(coll)
}

type DBCollections struct {
	usersCollection          *mongo.Collection
	otpsCollection           *mongo.Collection
	organizationsCollection  *mongo.Collection
	teamsCollection          *mongo.Collection
	activitiesCollection     *mongo.Collection
	datasCollections         *mongo.Collection
	uploadedFilesCollections *mongo.Collection
}

func (d *Database) GetAllCollections() *DBCollections {
	return &DBCollections{
		usersCollection:          d.GetCollection("users"),
		otpsCollection:           d.GetCollection("otps"),
		organizationsCollection:  d.GetCollection("organizations"),
		teamsCollection:          d.GetCollection("teams"),
		activitiesCollection:     d.GetCollection("activities"),
		datasCollections:         d.GetCollection("datas"),
		uploadedFilesCollections: d.GetCollection("uploaded_files"),
	}
}
