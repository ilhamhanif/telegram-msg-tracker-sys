package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
)

func (u *PubsubData) publishToPubSub() error {

	/*
		A method to publish the `raw` telegram update message
		Into BigQuery as logging.
	*/

	// Setup PubSub client.
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Setup the PubSub Message.
	jsonData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("publishToPubSub: Error marshalling struct: %w", err)
	}

	// Publish message to PubSub.
	t := client.Topic(PUBSUB_TOPIC_LOGGER)
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

func (u *PubsubData) getUpdateMessageID(o *IdentificationResult) error {

	/*
		A method to append Update ID.
	*/

	o.UpdateId = u.Message.ID

	return nil

}

func (u *PubsubData) getUpdateTime(o *IdentificationResult) error {

	/*
		A method to convert Update Epoch Time
		to Date and Datetime.
	*/

	var updateEpoch int

	if u.ChannelPost != nil {
		updateEpoch = u.ChannelPost.Date
	} else if u.ChannelPost == nil {
		updateEpoch = u.Message.Date
	} else {
		return fmt.Errorf("getUpdateTime: unknown telegram message `update` condition: %v", *u)
	}

	tz, _ := time.LoadLocation("Asia/Jakarta")
	t := time.Unix(int64(updateEpoch), 0).In(tz)
	o.UpdateEpoch = updateEpoch
	o.UpdateDate = t.Format("2006-01-02")
	o.UpdateDatetime = t.Format("2006-01-02T15:04:05")

	return nil

}

func (u *PubsubData) getUpdateType(o *IdentificationResult) error {

	/*
		A method to determine Update Type.
		Allowed type (currently):
		1. Text
		2. Photo
	*/

	var errReturn = fmt.Errorf("getUpdateType: unknown telegram message `update` condition: %v", *u)

	if u.ChannelPost != nil {

		// Channels do not have any abilities to receive bot commands.
		if u.ChannelPost.Text != "" {
			o.Result.Type = "TEXT"
			o.Result.Text = u.ChannelPost.Text
		} else if u.ChannelPost.Photo != nil {
			o.Result.Type = "PHOTO"
			o.Result.Photo = u.ChannelPost.Photo
		} else {
			return errReturn
		}
		o.Result.ChatId = u.ChannelPost.Chat.ID

	} else if u.ChannelPost == nil {

		// Determine traditional Type.
		if u.Message.Text != "" {
			o.Result.Type = "TEXT"
			o.Result.Text = u.Message.Text
		} else if u.Message.Photo != nil {
			o.Result.Type = "PHOTO"
			o.Result.Photo = u.Message.Photo
		} else {
			return errReturn
		}
		o.Result.ChatId = u.Message.Chat.ID

		// Update the Type if a bot_command is found.
		for _, v := range u.Message.Entities {
			if v.Type == "bot_command" {
				o.Result.Type = "BOT COMMAND"
			}
		}

	} else {
		return errReturn
	}

	return nil

}
