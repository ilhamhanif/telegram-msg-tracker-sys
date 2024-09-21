package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type Attributes struct {
	BucketID           string `json:"bucketId"`
	EventTime          string `json:"eventTime"`
	EventType          string `json:"eventType"`
	NotificationConfig string `json:"notificationConfig"`
	ObjectGeneration   string `json:"objectGeneration"`
	ObjectID           string `json:"objectId"`
	PayloadFormat      string `json:"payloadFormat"`
}

type PubsubMessage struct {
	Attributes        Attributes `json:"attributes"`
	Data              string     `json:"data"`
	MessageIDPascal   string     `json:"messageID"`
	MessageID         string     `json:"message_id"`
	PublishTimePascal string     `json:"publishTime"`
	PublishTime       string     `json:"publish_time"`
}

type PubsubSubscription struct {
	Message         *PubsubMessage `json:"message"`
	Subscription    string         `json:"subscription"`
	DeliveryAttempt int8           `json:"deliveryAttempt"`
}

const URL = "http://localhost:8080/GCSBucketNotifLogger"

var pubsubMessage = PubsubMessage{
	Attributes: Attributes{
		BucketID:           "test",
		EventTime:          "2022-08-12T23:22:36.901891Z",
		EventType:          "OBJECT_FINALIZE",
		NotificationConfig: "test",
		ObjectGeneration:   "16603467",
		ObjectID:           "itm/v",
		PayloadFormat:      "JSON_API_V1",
	},
	Data:        "eyJraW5kIjogInN0b3JhZ2Ujb2JqZWN0IiwgImlkIjogInRlc3QuY3N2LzE2NjAzNDY1NTY4Mjc1NjciLCAic2VsZkxpbmsiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vc3RvcmFnZS92MS9iL3Rlc3QuY3N2IiwgIm5hbWUiOiAidGVzdC5jc3YiLCAiYnVja2V0IjogInRlc3QiLCAiZ2VuZXJhdGlvbiI6ICIxNjYwMzQ2NTU2ODI3NTY3IiwgIm1ldGFnZW5lcmF0aW9uIjogIjEiLCAiY29udGVudFR5cGUiOiAidGV4dC9jc3YiLCAidGltZUNyZWF0ZWQiOiAiMjAyMi0wOC0xMlQyMzoyMjozNi45MDFaIiwgInVwZGF0ZWQiOiAiMjAyMi0wOC0xMlQyMzoyMjozNi45MDFaIiwgInN0b3JhZ2VDbGFzcyI6ICJNVUxUSV9SRUdJT05BTCIsICJ0aW1lU3RvcmFnZUNsYXNzVXBkYXRlZCI6ICIyMDIyLTA4LTEyVDIzOjIyOjM2LjkwMVoiLCAic2l6ZSI6ICIyMjAiLCAibWQ1SGFzaCI6ICJ1bUVkYW5NdERSbEpVZkdBWGNYemp3PT0iLCAibWVkaWFMaW5rIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL2Rvd25sb2FkL3N0b3JhZ2UvdjEvYi90ZXN0LmNzdj9nZW5lcmF0aW9uPTE2NjAzNDY1NTY4Mjc1NjcmYWx0PW1lZGlhIiwgIm1ldGFkYXRhIjogeyJnb29nLXJlc2VydmVkLWZpbGUtbXRpbWUiOiAiMTY1OTg3NTQ2MiJ9LCAiY3JjMzJjIjogIlBHQk5nUT09IiwgImV0YWciOiAiQ0srUDVmVzR3dmtDRUFFPSJ9",
	MessageID:   "5333919906745759",
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
