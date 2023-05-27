package storage

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"stockinos.com/api/models"
)

func (d *Database) CreateUserIfNotExists(ctx context.Context, phoneNumber string) error {
	var user models.User

	if err := d.DB.WithContext(ctx).First(&user, "phone_number = ?", phoneNumber).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{PhoneNumber: phoneNumber}
			result := d.DB.Create(&user)

			if result.Error != nil {
				return fmt.Errorf("error when creating user: %w", result.Error)
			}
		}
	}

	return nil
}
