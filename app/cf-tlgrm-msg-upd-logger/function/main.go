package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

type PubsubData models.Update

type BqRow struct {
	UpdateMessage PubsubData
}

const PROJECT_ID = "protean-quanta-434205-p5"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "telegram_msg_log_update"

func init() {
	functions.HTTP("TelegramMsgUpdateLogger", TelegramMsgUpdateLogger)
}

func TelegramMsgUpdateLogger(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubMessage PubsubSubscription
	var pubsubData PubsubData

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode PubSub Message.
	if err := pubsubMessage.decodePubSubData(&pubsubData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert data to Google BigQuery.
	var bqRows = BqRow{
		UpdateMessage: pubsubData,
	}
	if err := bqRows.insertBqRows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
