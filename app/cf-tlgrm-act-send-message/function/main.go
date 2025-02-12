package function

import (
	"encoding/json"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

type ApiResult struct {
	StatusCode int            `json:"status_code"`
	Message    models.Message `json:"message"`
}

const BOT_TOKEN = "7536185035:AAEGzJBD1iomeooHuRYpZtW81R-OyOECsBg"
const URL = "https://api.telegram.org/bot" + BOT_TOKEN + "/sendMessage"
const PROJECT_ID = "protean-quanta-434205-p5"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "telegram_act_log_send_message"

func init() {
	functions.HTTP("TelegramSendMessage", TelegramSendMessage)
}

func TelegramSendMessage(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubSubscription PubsubSubscription
	var pubsubData PubsubData
	var apiResult ApiResult

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

	// Send Message to Telegram.
	if err := pubsubData.sendMessage(&apiResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Store send message log result.
	var bqRows = BqRow{
		PubsubData: pubsubData,
		ApiResult:  apiResult,
	}
	if err := bqRows.insertBqRows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
