package main

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

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC = "tlgrm_act_send_message"

var botSendMessageParams = BotSendMessageParams{
	UpdateID:       12412423321,
	UpdateEpoch:    1725978621,
	UpdateDate:     "2024-09-10",
	UpdateDatetime: "2024-09-10T21:30:21",
	Params: bot.SendMessageParams{
		ChatID: -1002157107054,
		Text:   "Test",
	},
}

func (sm *BotSendMessageParams) publishToPubSub() error {

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC)
	jsonData, err := json.Marshal(sm)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error encoding JSON: %w", err)
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonData,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("publishToPubSub: Error publishing to PubSub: %w", err)
	}

	return nil

}

func main() {

	// Send a messages through Pub/Sub.
	if err := botSendMessageParams.publishToPubSub(); err != nil {
		fmt.Printf("Error: %s", err)
	}

}
