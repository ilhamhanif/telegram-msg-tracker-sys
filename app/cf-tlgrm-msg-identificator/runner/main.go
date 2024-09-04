package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Data        string `json:"data"`
	MessageId   string `json:"message_id"`
	PublishTime string `json:"publish_time"`
}

type PubsubSubscription struct {
	Message *PubsubMessage `json:"message"`
}

const URL = "http://localhost:8080/TelegramMsgIdentificator"

var telegramMsgUpdate = TelegramApiModelUpdate{
	ID: 2323514213,
	Message: &models.Message{
		ID: 23,
		From: &models.User{
			ID:           124523124,
			IsBot:        false,
			FirstName:    "Test",
			LastName:     "Test",
			Username:     "Test",
			LanguageCode: "en",
		},
		Chat: models.Chat{
			ID:    -5242314,
			Title: "Test",
			Type:  "group",
		},
		Date: 1231352213,
		// Text: "Test",
		Photo: []models.PhotoSize{
			{
				FileID:       "12ea-d12efeaddwe1221e",
				FileUniqueID: "asdfwqeaddas",
				FileSize:     783,
				Height:       60,
				Width:        90,
			},
			{
				FileID:       "12ea-d12efeaddwe1221e",
				FileUniqueID: "asdfwqeaddas",
				FileSize:     783,
				Height:       60,
				Width:        90,
			},
		},
	},
}

var pubsubMessage = PubsubMessage{
	MessageId:   "12453421242435123",
	PublishTime: "2022-08-12T23:22:36.971Z",
}

var pubsubSubscription = PubsubSubscription{
	Message: &pubsubMessage,
}

func main() {

	// Setup message in JSON
	// mimic-ing real GCP Pub/Sub HTTP push message
	messageJson, err := json.Marshal(telegramMsgUpdate)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	pubsubSubscription.Message.Data = base64.StdEncoding.EncodeToString(messageJson)
	payloadJson, err := json.Marshal(pubsubSubscription)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// Sent the data to local endpoint
	// using HTTP POST
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payloadJson))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer resp.Body.Close()

	// Print response and status code
	// given from the API
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, string(body))

}
