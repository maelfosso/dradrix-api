package storage

import (
	"context"
	"fmt"

	"stockinos.com/api/models"
)

func (d *Database) SaveWAMessages(ctx context.Context, messages []models.WhatsappMessage) error {
	d.DB.WithContext(ctx).Create(&messages)

	for _, message := range messages {
		// d.log(message.ID)
		fmt.Println("ID Message : ", message.ID)
	}

	return nil
}

func (d *Database) SaveWAStatus(ctx context.Context, statuses []models.WhatsappStatus) error {
	d.DB.WithContext(ctx).Create(&statuses)

	for _, message := range statuses {
		// d.log(message.ID)
		fmt.Println("ID Message : ", message.ID)
	}

	return nil
}
