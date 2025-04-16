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

type AddUploadedFileParams struct {
	UploadedBy primitive.ObjectID
	ActivityId primitive.ObjectID
	FileKey    string
}

func (q *Queries) AddUploadedFile(ctx context.Context, arg AddUploadedFileParams) (*models.UploadedFile, error) {
	var file models.UploadedFile = models.UploadedFile{
		UploadedBy: arg.UploadedBy,
		ActivityId: arg.ActivityId,
		FileKey:    arg.FileKey,

		UploadedAt: time.Now(),
	}

	_, err := q.uploadedFilesCollections.InsertOne(ctx, file)
	if err != nil {
		return nil, err
	} else {
		return &file, nil
	}
}

type GetAllUploadedFilesParams struct {
	UploadedBy primitive.ObjectID
	ActivityId primitive.ObjectID
}

func (q *Queries) GetAllUploadedFiles(ctx context.Context, arg GetAllUploadedFilesParams) ([]*models.UploadedFile, error) {
	var files []*models.UploadedFile

	filter := bson.M{
		"uploaded_by": arg.UploadedBy,
		"activity_id": arg.ActivityId,
	}
	cursor, err := q.uploadedFilesCollections.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

type RemoveUploadedFileParams struct {
	UploadedBy primitive.ObjectID
	ActivityId primitive.ObjectID
	FileKey    string
}

func (q *Queries) RemoveUploadedFile(ctx context.Context, arg RemoveUploadedFileParams) error {
	filter := bson.M{
		"uploaded_by": arg.UploadedBy,
		"activity_id": arg.ActivityId,
		"file_key":    arg.FileKey,
	}

	_, err := q.uploadedFilesCollections.DeleteOne(
		ctx,
		filter,
	)
	return err
}

type RemoveAllUploadedFileParams struct {
	UploadedBy primitive.ObjectID
	ActivityId primitive.ObjectID
}

func (q *Queries) RemoveAllUploadedFile(ctx context.Context, arg RemoveUploadedFileParams) error {
	filter := bson.M{
		"user_id":     arg.UploadedBy,
		"activity_id": arg.ActivityId,
	}

	_, err := q.uploadedFilesCollections.DeleteMany(
		ctx,
		filter,
	)
	return err
}
