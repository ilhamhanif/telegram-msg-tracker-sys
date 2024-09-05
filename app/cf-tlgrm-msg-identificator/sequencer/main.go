package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

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

func publishToPubSub(message []byte) error {

	// Setup PubSub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC)
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

func main() {

	// Send 100 messages through Pub/Sub
	for i := 0; i <= 100; i++ {

		telegramMsgUpdate.ID = 123523412 + int64(i)
		jsonData, err := json.Marshal(telegramMsgUpdate)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		if err := publishToPubSub(jsonData); err != nil {
			fmt.Printf("Error: %s", err)
		}
		fmt.Printf("%d %d\n", i, telegramMsgUpdate.ID)

	}

}
