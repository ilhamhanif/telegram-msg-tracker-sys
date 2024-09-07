package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC = "tlgrm_msg_identificator"

var telegramMsgUpdate = TelegramApiModelUpdate{
	Message: &models.Message{
		ID: 28,
		From: &models.User{
			ID:           123123412,
			IsBot:        false,
			FirstName:    "Test",
			LastName:     "Test",
			Username:     "Test",
			LanguageCode: "en",
		},
		Chat: models.Chat{
			ID:    -42432431,
			Title: "Test",
			Type:  "group",
		},
		Date: 674763452,
		Text: "Test",
	},
}

func (u *TelegramApiModelUpdate) publishToPubSub() error {

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC)
	jsonData, err := json.Marshal(u)
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

	// Send 100 messages through Pub/Sub.
	for i := 0; i <= 100; i++ {

		telegramMsgUpdate.ID = 123523412 + int64(i)
		if err := telegramMsgUpdate.publishToPubSub(); err != nil {
			fmt.Printf("Error: %s", err)
		}
		fmt.Printf("%d %d\n", i, telegramMsgUpdate.ID)

	}

}
