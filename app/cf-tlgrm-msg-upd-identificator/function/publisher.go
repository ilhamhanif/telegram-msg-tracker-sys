package function

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/go-telegram/bot"
)

type BotSendMessageParams struct {
	UpdateID       int                   `json:"update_id"`
	UpdateEpoch    int                   `json:"update_epoch"`
	UpdateDate     string                `json:"update_date"`
	UpdateDatetime string                `json:"update_datetime"`
	Params         bot.SendMessageParams `json:"params"`
}

func (sm *BotSendMessageParams) sendMessage() error {

	// Setup PubSub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("sendMessage: Error initiating PubSub client: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC_SEND_MESSAGE)
	jsonData, err := json.Marshal(sm)
	if err != nil {
		return fmt.Errorf("sendMessage: Error encoding JSON: %w", err)
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonData,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("sendMessage: Error publishing to PubSub: %w", err)
	}

	return nil

}
