package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC_REGULATOR = "telegram_msg_regulator"
const PUBSUB_TOPIC_LOGGER = "telegram_msg_update_logger"

var telegramMsgUpdate TelegramApiModelUpdate

func init() {
	functions.HTTP("TelegramMsgUpdateForwarder", TelegramMsgUpdateForwarder)
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

func TelegramMsgUpdateForwarder(w http.ResponseWriter, r *http.Request) {

	// Receive and parse HTTP push data message
	// from Telegram Webhook
	if err := json.NewDecoder(r.Body).Decode(&telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(telegramMsgUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// // Forward the message to REGULATOR through Pub/Sub
	// if err := publishToPubSub(PUBSUB_TOPIC_REGULATOR, jsonData); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Forward the message to LOGGER through Pub/Sub
	if err := publishToPubSub(PUBSUB_TOPIC_LOGGER, jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`
	fmt.Fprint(w, "ok")

}
