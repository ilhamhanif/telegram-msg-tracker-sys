package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type BotFileDownloadParams struct {
	UpdateID       int                 `json:"update_id"`
	UpdateEpoch    int                 `json:"update_epoch"`
	UpdateDate     string              `json:"update_date"`
	UpdateDatetime string              `json:"update_datetime"`
	Files          []map[string]string `json:"files"`
}

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC = "tlgrm_utils_file_downloader"

var botFileDownloadParams = BotFileDownloadParams{
	UpdateID:       53726587267,
	UpdateEpoch:    1725723760,
	UpdateDate:     "2024-09-07",
	UpdateDatetime: "2024-09-07T22:42:40",
	Files: []map[string]string{
		{
			"file_id":        "AgACAgUAAyEFAASAktduAAMOZuBX_cAckrjv39BaJsVqBqu7aUIAAsi_MRs5ggFXIl4z2xzgmUIBAAMCAANzAAM2BA",
			"file_unique_id": "AQADyL8xGzmCAVd4",
		},
	},
}

func (sm *BotFileDownloadParams) publishToPubSub() error {

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
	if err := botFileDownloadParams.publishToPubSub(); err != nil {
		fmt.Printf("Error: %s", err)
	}

}
