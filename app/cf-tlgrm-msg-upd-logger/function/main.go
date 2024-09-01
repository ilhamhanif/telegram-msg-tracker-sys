package function

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
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

type BqRow struct {
	UpdateMessage TelegramApiModelUpdate
}

const PROJECT_ID = "protean-quanta-434205-p5"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "telegram_msg_log_update"

func init() {
	functions.HTTP("TelegramMsgUpdateLogger", TelegramMsgUpdateLogger)
}

func (r *BqRow) Save() (map[string]bigquery.Value, string, error) {

	/*
		A method to format BigQuery row record
	*/

	tz, _ := time.LoadLocation("Asia/Jakarta")
	currDatetime := time.Now().In(tz)
	logDatetime := currDatetime.Format("2006-01-02T15:04:05")
	logEpoch := currDatetime.Format("20060102150405")
	logDate := currDatetime.Format("2006-01-02")
	updateId := r.UpdateMessage.ID
	update, _ := json.Marshal(r.UpdateMessage)
	updateStr := string(update)

	return map[string]bigquery.Value{
		"update_id":    updateId,
		"update":       updateStr,
		"log_date":     logDate,
		"log_datetime": logDatetime,
		"log_epoch":    logEpoch,
	}, bigquery.NoDedupeID, nil
}

func insertBqRows(rows []*BqRow) error {

	/*
		A function to insert row records to GCP BigQuery
	*/

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return err
	}
	defer client.Close()

	inserter := client.Dataset(BQ_DATASET_NAME).Table(BQ_TABLE_NAME).Inserter()
	if err := inserter.Put(ctx, rows); err != nil {
		return err
	}

	return nil
}

func TelegramMsgUpdateLogger(w http.ResponseWriter, r *http.Request) {

	var pubsubMessage PubsubSubscription
	var telegramMsgUpdate TelegramApiModelUpdate

	// Receive and parse GCP Pub/Sub HTTP push data message
	// and Decode the data
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pubsubMessageDataDecoded, _ := base64.StdEncoding.DecodeString(pubsubMessage.Message.Data)

	// Convert string JSON
	if err := json.Unmarshal(pubsubMessageDataDecoded, &telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert data to Google BigQuery
	bqRows := []*BqRow{
		{UpdateMessage: telegramMsgUpdate},
	}
	if err := insertBqRows(bqRows); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`
	fmt.Fprint(w, "ok")

}
