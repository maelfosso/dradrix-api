package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetUserByUsernameParams struct {
	PhoneNumber string
}

func (q *Queries) GetUserByUsername(ctx context.Context, arg GetUserByUsernameParams) error {
	var result bson.M

	err := q.usersCollection.FindOne(
		ctx,
		bson.D{{Key: "phoneNumber", Value: arg.PhoneNumber}},
	).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
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
