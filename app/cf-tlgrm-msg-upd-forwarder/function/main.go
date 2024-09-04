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

func init() {
	functions.HTTP("TelegramMsgUpdateForwarder", TelegramMsgUpdateForwarder)
}

func (u *TelegramApiModelUpdate) publishToPubSub() error {

	// Setup PubSub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC_IDENTIFICATOR)
	jsonData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error marshalling struct: %w", err)
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonData,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("publishToPubSub: result.Get: %w", err)
	}

	return nil

}

func TelegramMsgUpdateForwarder(w http.ResponseWriter, r *http.Request) {

	var telegramMsgUpdate TelegramApiModelUpdate

	// Receive and parse HTTP push data message
	// from Telegram Webhook
	if err := json.NewDecoder(r.Body).Decode(&telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the message to IDENTIFICATOR through Pub/Sub
	if err := telegramMsgUpdate.publishToPubSub(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`
	fmt.Fprint(w, "ok")

}
