package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var sendMessageTextJsonTemplate = `{
	"messaging_product": "whatsapp",
	"recipient_type": "individual",
	"to": "%s",
	"type": "text",
	"text": {
			"preview_url": false,
			"body": "%s"
	}
}`

func SendMessageText(to, message string) (*WhatsappSendTextMessageResponse, error) {
	jsonBody := []byte(fmt.Sprintf(sendMessageTextJsonTemplate, to, message))
	bodyReader := bytes.NewReader(jsonBody)
	log.Println("send message text data : ", string(jsonBody))

	requestUrl := fmt.Sprintf(
		"https://graph.facebook.com/%s/%s/messages",
		os.Getenv("WHATSAPP_API_VERSION"), os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
	)
	log.Println("send message text request url: ", requestUrl)

	req, err := http.NewRequest(
		http.MethodPost,
		requestUrl,
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WHATSAPP_USER_ACCESS_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: error making http request: %s", err)
	}

	log.Println("Client: got response!")
	log.Println("client: status code: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s", err)
	}

	log.Println("client: response body: %s", string(resBody))

	var data WhatsappSendTextMessageResponse
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshilling response body : %s", err)
	}

	return &data, nil
}
