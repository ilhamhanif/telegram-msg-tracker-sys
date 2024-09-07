package function

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func (u *TelegramApiModelUpdate) publishToPubSub(pubsub_topic string) error {

	/*
		A method to publish the message to PubSub.
	*/

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error initiating PubSub client: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub.
	t := client.Topic(pubsub_topic)
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
