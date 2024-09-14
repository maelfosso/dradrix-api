package storage

import (
	"context"
	"errors"
	"fmt"
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
	FirstName   string
	LastName    string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*models.User, error) {
	var user models.User = models.User{
		Id:          primitive.NewObjectID(),
		PhoneNumber: arg.PhoneNumber,
		FirstName:   arg.FirstName,
		LastName:    arg.LastName,

		Preferences: models.UserPreferences{
			CurrentStatus: "is-creating-account",
		},
	}

	result, err := q.usersCollection.InsertOne(ctx, user)
	log.Println("[CreateUser] ", result.InsertedID, user.Id)
	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

type UpdateUserProfileParams struct {
	Id        primitive.ObjectID
	FirstName string
	LastName  string
	Email     string
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (*models.User, error) {
	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": bson.M{
			"first_name": arg.FirstName,
			"last_name":  arg.LastName,
			"email":      arg.Email,
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

type UpdateUserPreferencesParams struct {
	Id primitive.ObjectID

	Changes map[string]any
}

func (q *Queries) UpdateUserPreferences(ctx context.Context, arg UpdateUserPreferencesParams) (*models.User, error) {
	set := bson.M{}
	for name, value := range arg.Changes {
		field := fmt.Sprintf("preferences.%s", name)
		set[field] = value
	}

	filter := bson.M{
		"_id": arg.Id,
	}
	update := bson.M{
		"$set": set,
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
