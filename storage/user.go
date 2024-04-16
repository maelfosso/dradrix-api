package storage

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"gorm.io/gorm"
// 	"stockinos.com/api/models"
// )

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
