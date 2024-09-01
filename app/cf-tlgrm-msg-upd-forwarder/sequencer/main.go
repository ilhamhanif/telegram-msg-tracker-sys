package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

const URL = "https://cf-tlgrm-msg-upd-forwarder-j4tflaldfq-et.a.run.app"

var telegramMsgUpdate = TelegramApiModelUpdate{
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
		Text: "Test",
	},
}

func main() {

	for i := 0; i <= 10000; i++ {
		// Setup message in JSON
		// mimic-ing real GCP Pub/Sub HTTP push message
		telegramMsgUpdate.Message.ID = 2131241246 + i
		messageJson, err := json.Marshal(telegramMsgUpdate)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		// Sent the data to local endpoint
		// using HTTP POST
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(messageJson))
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

}
