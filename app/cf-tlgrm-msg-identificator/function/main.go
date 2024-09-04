package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

// https://cloud.google.com/pubsub/docs/push#properties_of_a_push_subscription
type PubsubMessage struct {
	Attributes        map[string]string `json:"attributes"`
	Data              string            `json:"data"`
	MessageIdPascal   string            `json:"messageId"`
	MessageId         string            `json:"message_id"`
	PublishTimePascal string            `json:"publishTime"`
	PublishTime       string            `json:"publish_time"`
}

type PubsubSubscription struct {
	Message         *PubsubMessage `json:"message"`
	Subscription    string         `json:"subscription"`
	DeliveryAttempt int8           `json:"deliveryAttempt"`
}

type TelegramApiModelUpdate models.Update
type IdentificationResult struct {
	UpdateId       int
	UpdateEpoch    int
	UpdateDate     string
	UpdateDatetime string
	Result         struct {
		Type  string
		Text  string
		Photo []models.PhotoSize
	}
}

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC_LOGGER = "telegram_msg_update_logger"
const PUBSUB_TOPIC_SEND_MESSAGE = "telegram_msg_action_send_message"

func (pm *PubsubSubscription) decodePubSubData(v *TelegramApiModelUpdate) error {

	/*
		A method to Decode a PubSub data to v
	*/

	// Convert (decode) string JSON
	pubsubMessageDataDecoded, _ := base64.StdEncoding.DecodeString(pm.Message.Data)
	if err := json.Unmarshal(pubsubMessageDataDecoded, &v); err != nil {
		return fmt.Errorf("decodePubSubData: Error decoding PubSub data: %w", err)
	}

	return nil

}

func (u *TelegramApiModelUpdate) publishRawDataToPubSub() error {

	// Setup PubSub client
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("publishRawDataToPubSub: Error creating NewClient: %w", err)
	}
	defer client.Close()

	// Setup the PubSub Message
	pubsubMessageDataDecoded, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("publishRawDataToPubSub: Error marshalling struct: %w", err)
	}

	// Publish message to PubSub
	t := client.Topic(PUBSUB_TOPIC_LOGGER)
	result := t.Publish(ctx, &pubsub.Message{
		Data: pubsubMessageDataDecoded,
	})

	// Block until the result is returned
	// and a server-generated ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("publishRawDataToPubSub: Error result.Get: %w", err)
	}
	return nil
}

func (u *TelegramApiModelUpdate) getUpdateMessageID(o *IdentificationResult) error {

	/*
		A method to append Update ID.
	*/

	o.UpdateId = u.Message.ID

	return nil
}

func (u *TelegramApiModelUpdate) getUpdateTime(o *IdentificationResult) error {

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

func (u *TelegramApiModelUpdate) getUpdateType(o *IdentificationResult) error {

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

func init() {
	functions.HTTP("TelegramMsgIdentificator", TelegramMsgIdentificator)
}

func TelegramMsgIdentificator(w http.ResponseWriter, r *http.Request) {

	var pubsubMessage PubsubSubscription
	var telegramMsgUpdate TelegramApiModelUpdate
	var identificationResult IdentificationResult

	// Receive and parse GCP Pub/Sub HTTP push data message
	// and Decode the data from HTTP Binary
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode PubSub data to get the `raw` data.
	if err := pubsubMessage.decodePubSubData(&telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the `raw` message to LOGGER through Pub/Sub
	if err := telegramMsgUpdate.publishRawDataToPubSub(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get update characteristics the result
	if err := telegramMsgUpdate.getUpdateMessageID(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := telegramMsgUpdate.getUpdateType(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := telegramMsgUpdate.getUpdateTime(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(identificationResult)

	// Return `ok`
	fmt.Fprint(w, "ok")

}
