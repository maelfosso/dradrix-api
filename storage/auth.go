package storage

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"

// 	"gorm.io/gorm"
// 	"stockinos.com/api/models"
// )

// const OTP_VALIDITY_MINUTES = 5

// func (d *Database) CreateOTP(ctx context.Context, otp models.OTP) error {
// 	var m models.OTP

// 	// Check if an instance with the phone number that is active, already exists
// 	err := d.DB.WithContext(ctx).First(&m, "phone_number = ? AND active = ?", otp.PhoneNumber, true).Error
// 	if err == nil {
// 		// If so, invalidate it
// 		m.Active = false
// 		d.DB.WithContext(ctx).Save(&m)
// 	} else {
// 		if !errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("error when looking for existing active OTP: %w", err)
// 		}
// 	}

// 	// Save the current one
// 	d.DB.WithContext(ctx).Create(&otp)
// 	return nil
// }

// func (d *Database) CheckOTP(ctx context.Context, phoneNumber, pinCode string) (*models.OTP, error) {
// 	var otp models.OTP

// 	wQuery := d.DB.WithContext(ctx).Where("phone_number = ? AND active = ?", phoneNumber, true).First(&otp)
// 	if err := wQuery.Error; errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, errors.New("OTP_NOT_EXISTS_WITH_PHONE_NUMBER")
// 	}

// 	if !otp.Active {
// 		return nil, errors.New("OTP_NO_LONGER_ACTIVE")
// 	}

// 	if otp.PinCode != pinCode {
// 		return nil, errors.New("OTP_WRONG")
// 	}

// 	ellapsed := time.Since(otp.CreatedAt)
// 	if ellapsed.Minutes() >= OTP_VALIDITY_MINUTES {
// 		return nil, errors.New("OTP_ALREADY_EXPIRED")
// 	}

// 	return &otp, nil
// }

// func (d *Database) SaveOTP(ctx context.Context, otp models.OTP) error {
// 	d.DB.WithContext(ctx).Save(&otp)
// 	return nil
// }
