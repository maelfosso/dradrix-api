package storage

import (
	"context"

	"stockinos.com/api/models"
)

func (d *Database) SavePinCode(ctx context.Context, pinCode models.PinCodeOTP) error {
	d.DB.WithContext(ctx).Create(&pinCode)
	return nil
}
