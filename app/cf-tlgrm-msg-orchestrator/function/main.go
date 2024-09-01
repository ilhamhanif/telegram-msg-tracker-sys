package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Attributes        map[string]string `json:"attributes"`
	Data              string            `json:"data"`
	MessageIdPascal   string            `json:"messageId"`
	MessageId         string            `json:"message_id"`
	PublishTimePascal string            `json:"publishTime"`
	PublishTime       string            `json:"publish_time"`
}

type PubsubSubscription struct {
	Message         *PubsubMessage `json:"message"`
	Subscription    string         `json:"subscription"`
	DeliveryAttempt int8           `json:"deliveryAttempt"`
}

type TelegramApiModelUpdate models.Update

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC_LOGGER = "telegram_msg_update_logger"

func init() {
	functions.HTTP("TelegramMsgOrchestrator", TelegramMsgOrchestrator)
}

func publishToPubSub(pubsub_topic string, message []byte) error {

	// Setup PubSub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(pubsub_topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	return nil

}

func TelegramMsgOrchestrator(w http.ResponseWriter, r *http.Request) {

	var pubsubMessage PubsubSubscription
	var telegramMsgUpdate TelegramApiModelUpdate

	// Receive and parse GCP Pub/Sub HTTP push data message
	// and Decode the data
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubsubMessageDataDecoded, _ := base64.StdEncoding.DecodeString(pubsubMessage.Message.Data)

	// Convert string JSON
	if err := json.Unmarshal(pubsubMessageDataDecoded, &telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check message with the following condition.
	//
	// 1. If Text:
	// 2. If File (Photo, Video, ... ):
	//
	// For all: Forward the message to LOGGER through Pub/Sub

	// Forward the message to LOGGER through Pub/Sub
	if err := publishToPubSub(PUBSUB_TOPIC_LOGGER, pubsubMessageDataDecoded); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`
	fmt.Fprint(w, "ok")

}
