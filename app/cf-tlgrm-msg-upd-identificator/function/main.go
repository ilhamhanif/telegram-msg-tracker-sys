package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC_SEND_MESSAGE = "tlgrm_act_send_message"
const PUBSUB_TOPIC_FILE_DOWNLOADER = "tlgrm_utils_file_downloader"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "telegram_msg_log_identification"

func init() {
	functions.HTTP("TelegramMsgUpdIdentificator", TelegramMsgUpdIdentificator)
}

func TelegramMsgUpdIdentificator(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubSubscription PubsubSubscription
	var pubsubData PubsubData
	var identificationResult IdentificationResult

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubSubscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode PubSub data to get the `raw` data.
	if err := pubsubSubscription.decodePubSubData(&pubsubData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Identify the message.
	if err := pubsubData.getUpdateMessageID(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := pubsubData.getUpdateType(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := pubsubData.getUpdateTime(&identificationResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle the identification result.
	if err := identificationResult.check(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store identification result.
	var bqRows = BqRow{
		IdentificationResult: identificationResult,
	}
	if err := bqRows.insertBqRows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
