package storage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"stockinos.com/api/models"
)

type DoesUserExistsParams struct {
	PhoneNumber string
}

func (q *Queries) DoesUserExists(ctx context.Context, arg DoesUserExistsParams) (*models.User, error) {
	var user models.User

	err := q.usersCollection.FindOne(
		ctx,
		bson.D{{Key: "phoneNumber", Value: arg.PhoneNumber}},
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
		bson.D{{Key: "phoneNumber", Value: arg.PhoneNumber}},
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

// func (d *Database) CreateUserIfNotExists(ctx context.Context, phoneNumber string) error {
// 	var user models.User

// 	if err := d.DB.WithContext(ctx).First(&user, "phone_number = ?", phoneNumber).Error; err != nil { // errors.Is(err, gorm.ErrRecordNotFound) {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			user = models.User{PhoneNumber: phoneNumber}
// 			result := d.DB.Create(&user)

// 			return result.Error
// 			// if result.Error != nil {
// 			// 	return fmt.Errorf("error when creating user: %w", result.Error)
// 			// } else
// 		} else {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (d *Database) FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
// 	var user models.User

// 	if err := d.DB.WithContext(ctx).First(&user, "phone_number = ?", phoneNumber).Error; err != nil { // } errors.Is(err, gorm.ErrRecordNotFound) {
// 		// return nil, fmt.Errorf("error no user found")
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, fmt.Errorf("error no user found")
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	return &user, nil
// }
