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

func getWhatsappRequestURL(url string) string {
	return fmt.Sprintf(
		url,
		os.Getenv("WHATSAPP_API_VERSION"), os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
	)
}

var jsonForSendMessageText = `{
	"messaging_product": "whatsapp",
	"recipient_type": "individual",
	"to": "%s",
	"type": "text",
	"text": {
			"preview_url": false,
			"body": "%s"
	}
}`

func SendMessageText(to, message string) (*WhatsappSendMessageResponse, error) {
	jsonBody := []byte(fmt.Sprintf(jsonForSendMessageText, to, message))
	bodyReader := bytes.NewReader(jsonBody)
	log.Println("send message text data : ", string(jsonBody))

	requestUrl := getWhatsappRequestURL("https://graph.facebook.com/%s/%s/messages")
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
	log.Printf("client: status code: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %s", err)
	}

	log.Printf("client: response body: %s", string(resBody))

	var data WhatsappSendMessageResponse
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshilling response body : %s", err)
	}

	return &data, nil
}

var jsonForSendMessageTemplateText = `{
	"messaging_product": "whatsapp",
	"recipient_type": "individual",
	"to": "%s",
	"type": "template",
	"template": {
		"name": "%s",
		"language": {
			"code": "%s"
		},
		"components": [
			{
				"type": "body",
				"parameters": %s
			}
		]
	}
}`

// [
// 	{
// 		"type": "text",
// 		"text": "text-string"
// 	},
// 	{
// 		"type": "currency",
// 		"currency": {
// 			"fallback_value": "$100.99",
// 			"code": "USD",
// 			"amount_1000": 100990
// 		}
// 	},
// 	{
// 		"type": "date_time",
// 		"date_time": {
// 			"fallback_value": "February 25, 1977",
// 			"day_of_week": 5,
// 			"year": 1977,
// 			"month": 2,
// 			"day_of_month": 25,
// 			"hour": 15,
// 			"minute": 33,
// 			"calendar": "GREGORIAN"
// 		}
// 	}
// ]

func SendMessageTextFromTemplate(to, template, language, parameters string) (*WhatsappSendMessageResponse, error) {
	jsonBody := []byte(
		fmt.Sprintf(
			jsonForSendMessageTemplateText, to, template, language, parameters,
		),
	)
	bodyReader := bytes.NewReader(jsonBody)

	requestUrl := getWhatsappRequestURL("https://graph.facebook.com/%s/%s/messages")

	req, err := http.NewRequest(
		http.MethodPost,
		requestUrl,
		bodyReader,
	)
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("WHATSAPP_USER_ACCESS_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %w", err)
	}

	log.Println("client: got response!")
	log.Printf("client: status code: %d", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("client: could not read response body: %w", err)
	}

	log.Printf("client: response body: %s", string(resBody))

	var data WhatsappSendMessageResponse
	err = json.Unmarshal(resBody, &data)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshelling response body: %w", err)
	}

	return &data, nil
}

func SendWoZOTP(to, language, pinCode string) (string, error) {
	parameters := fmt.Sprintf(`[
		{
			"type": "text",
			"text": "%s"
		}
	]`, pinCode)
	template := "woz_otp"

	res, err := SendMessageTextFromTemplate(to, template, language, parameters)
	if err != nil {
		return "", err
	}
	return res.Messages[0].ID, nil
}
