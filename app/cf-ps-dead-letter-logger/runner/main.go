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

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type Attributes struct {
	CloudPubSubDeadLetterSourceDeliveryCount       string
	CloudPubSubDeadLetterSourceSubscription        string
	CloudPubSubDeadLetterSourceSubscriptionProject string
	CloudPubSubDeadLetterSourceTopicPublishTime    string
}

type PubsubMessage struct {
	Attributes  Attributes
	Data        string `json:"data"`
	MessageID   string `json:"message_id"`
	PublishTime string `json:"publish_time"`
}

type PubsubSubscription struct {
	Message *PubsubMessage `json:"message"`
}

const URL = "http://localhost:8080/PubSubDeadLetterLogger"

var attributes = Attributes{
	CloudPubSubDeadLetterSourceDeliveryCount:       "5",
	CloudPubSubDeadLetterSourceSubscription:        "testTest",
	CloudPubSubDeadLetterSourceSubscriptionProject: "testTest",
	CloudPubSubDeadLetterSourceTopicPublishTime:    "2023-06-28T07:36:09.478+00:00",
}

var pubsubMessage = PubsubMessage{
	Attributes:  attributes,
	Data:        "eyJ0ZXN0IjoidGVzdCJ9",
	MessageID:   "12453421242435123",
	PublishTime: "2022-08-12T23:22:36.971Z",
}

var pubsubSubscription = PubsubSubscription{
	Message: &pubsubMessage,
}

func main() {

	// Setup message in JSON
	// mimic-ing real GCP Pub/Sub HTTP push message.
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
