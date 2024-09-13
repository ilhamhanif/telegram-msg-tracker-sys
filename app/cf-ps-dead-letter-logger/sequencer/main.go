package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Attributes struct {
	CloudPubSubDeadLetterSourceDeliveryCount       string
	CloudPubSubDeadLetterSourceSubscription        string
	CloudPubSubDeadLetterSourceSubscriptionProject string
	CloudPubSubDeadLetterSourceTopicPublishTime    string
}

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC = "pubsub_log_dead_letter"

var attributes = map[string]string{
	"CloudPubSubDeadLetterSourceDeliveryCount":       "5",
	"CloudPubSubDeadLetterSourceSubscription":        "testTest",
	"CloudPubSubDeadLetterSourceSubscriptionProject": "testTest",
	"CloudPubSubDeadLetterSourceTopicPublishTime":    "2023-06-28T07:36:09.478+00:00",
}

func publishToPubSub() error {

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC)
	result := t.Publish(ctx, &pubsub.Message{
		Data:       []byte("eyJ0ZXN0IjoidGVzdCJ9"),
		Attributes: attributes,
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

		if err := publishToPubSub(); err != nil {
			fmt.Printf("Error: %s", err)
		}
		fmt.Printf("%d\n", i)
	}

}
