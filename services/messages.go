package services

import (
	"context"
	"time"

	"stockinos.com/api/models"
	"stockinos.com/api/requests"
	"stockinos.com/api/storage"
)

func HandleMessageSentByWoZ(db storage.Database, data models.WhatsAppMessage) error { // (*requests.WhatsAppSendTextMessageResponse, error) {
	response, err := requests.SendMessageText(data.To, data.Text.Body)
	if err != nil {
		return err
	}
	data.ID = response.Messages[0].ID

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = db.SaveWAMessages(ctx, []models.WhatsAppMessage{data})
	if err != nil {
		return err
	}
	return nil
}
