package services

import (
	"stockinos.com/api/requests"
)

func WASendTextMessage(to, body string) (string, error) { // (*requests.WhatsappSendTextMessageResponse, error) {
	response, err := requests.SendMessageText(to, body)
	if err != nil {
		return "", err
	}
	msgId := response.Messages[0].ID

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// err = db.SaveWAMessages(ctx, []models.WhatsappMessage{data})
	// if err != nil {
	// 	return err
	// }
	return msgId, nil
}
