package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Data        string `json:"data"`
	MessageID   string `json:"message_id"`
	PublishTime string `json:"publish_time"`
}

type PubsubSubscription struct {
	Message *PubsubMessage `json:"message"`
}

type BotFileDownloadParams struct {
	UpdateID       int                 `json:"update_id"`
	UpdateEpoch    int                 `json:"update_epoch"`
	UpdateDate     string              `json:"update_date"`
	UpdateDatetime string              `json:"update_datetime"`
	Files          []map[string]string `json:"files"`
}

const URL = "http://localhost:8080/TelegramUtilsFileDownloader"

var botFileDownloadParams = BotFileDownloadParams{
	UpdateID:       53726587267,
	UpdateEpoch:    1725723760,
	UpdateDate:     "2024-09-07",
	UpdateDatetime: "2024-09-07T22:42:40",
	Files: []map[string]string{
		{
			"file_id":        "AgACAgUAAyEFAASAktduAAMOZuBX_cAckrjv39BaJsVqBqu7aUIAAsi_MRs5ggFXIl4z2xzgmUIBAAMCAANzAAM2BA",
			"file_unique_id": "AQADyL8xGzmCAVd4",
		},
	},
}

var pubsubMessage = PubsubMessage{
	MessageID:   "12453421242435123",
	PublishTime: "2022-08-12T23:22:36.971Z",
}

var pubsubSubscription = PubsubSubscription{
	Message: &pubsubMessage,
}

func main() {

	// Setup message in JSON
	// mimic-ing real GCP Pub/Sub HTTP push message.
	messageJson, err := json.Marshal(botFileDownloadParams)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	pubsubSubscription.Message.Data = base64.StdEncoding.EncodeToString(messageJson)
	payloadJson, err := json.Marshal(pubsubSubscription)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// Sent the data to local endpoint
	// using HTTP POST.
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
	// given from the API.
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, string(body))

}
