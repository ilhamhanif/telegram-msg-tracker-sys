package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Attributes struct {
		CloudPubSubDeadLetterSourceDeliveryCount       string
		CloudPubSubDeadLetterSourceSubscription        string
		CloudPubSubDeadLetterSourceSubscriptionProject string
		CloudPubSubDeadLetterSourceTopicPublishTime    string
	}
	Data              string `json:"data"`
	MessageIDPascal   string `json:"messageID"`
	MessageID         string `json:"message_id"`
	PublishTimePascal string `json:"publishTime"`
	PublishTime       string `json:"publish_time"`
}

type PubsubSubscription struct {
	Message         *PubsubMessage `json:"message"`
	Subscription    string         `json:"subscription"`
	DeliveryAttempt int8           `json:"deliveryAttempt"`
}

const PROJECT_ID = "protean-quanta-434205-p5"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "pubsub_log_dead_letter"

func init() {
	functions.HTTP("PubSubDeadLetterLogger", PubSubDeadLetterLogger)
}

func PubSubDeadLetterLogger(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubSubscription PubsubSubscription

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubSubscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert data to Google BigQuery.
	var bqRows = BqRow{
		PubsubSubscription: pubsubSubscription,
	}
	if err := bqRows.insertBqRows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
