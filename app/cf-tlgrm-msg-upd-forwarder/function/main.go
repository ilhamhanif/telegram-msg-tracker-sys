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
const PUBSUB_TOPIC_ORCHESTRATOR = "telegram_msg_orchestrator"

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

	var telegramMsgUpdate TelegramApiModelUpdate

	// Receive and parse HTTP push data message
	// from Telegram Webhook
	fmt.Println("---")
	fmt.Println(telegramMsgUpdate)
	if err := json.NewDecoder(r.Body).Decode(&telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(telegramMsgUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(jsonData))

	// Forward the message to ORCHESTRATOR through Pub/Sub
	if err := publishToPubSub(PUBSUB_TOPIC_ORCHESTRATOR, jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`
	fmt.Fprint(w, "ok")

}
