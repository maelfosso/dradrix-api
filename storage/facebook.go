package storage

import (
	"context"
	"fmt"

	"stockinos.com/api/models"
)

func (d *Database) SaveWAMessages(ctx context.Context, messages []models.WhatsAppMessage) error {
	d.DB.WithContext(ctx).Create(&messages)

	for _, message := range messages {
		// d.log(message.ID)
		fmt.Println("ID Message : ", message.ID)
	}

	return nil
}
