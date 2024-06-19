package storage

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"stockinos.com/api/models"
)

type DoesUserExistsParams struct {
	PhoneNumber string
}

func (q *Queries) DoesUserExists(ctx context.Context, arg DoesUserExistsParams) (*models.User, error) {
	var user models.User

	err := q.usersCollection.FindOne(
		ctx,
		bson.D{{Key: "phone_number", Value: arg.PhoneNumber}},
	).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

type GetUserByPhoneNumberParams struct {
	PhoneNumber string
}

func (q *Queries) GetUserByPhoneNumber(ctx context.Context, arg GetUserByPhoneNumberParams) (*models.User, error) {
	var user models.User

	err := q.usersCollection.FindOne(
		ctx,
		bson.D{{Key: "phone_number", Value: arg.PhoneNumber}},
	).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}

type CreateUserParams struct {
	PhoneNumber string
	Name        string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error) {
	var user models.User = models.User{
		Id:          primitive.NewObjectID(),
		PhoneNumber: arg.PhoneNumber,
		Name:        arg.Name,
	}

	result, err := q.usersCollection.InsertOne(ctx, user)
	log.Println("[CreateUser] ", result.InsertedID, user.Id)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

type UpdateUserNameParams struct {
	Id        primitive.ObjectID
	FirstName string
	LastName  string
}

func (q *Queries) UpdateUserName(ctx context.Context, arg UpdateUserNameParams) (*models.User, error) {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"first_name": arg.FirstName,
			"last_name":  arg.LastName,
		},
	}
	after := options.After

	var user models.User
	err := q.usersCollection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		&options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		},
	).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &user, nil
}
