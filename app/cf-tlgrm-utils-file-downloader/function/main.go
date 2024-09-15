package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const PROJECT_ID = "protean-quanta-434205-p5"
const BOT_TOKEN = "7536185035:AAEGzJBD1iomeooHuRYpZtW81R-OyOECsBg"
const BASE_URL = "https://api.telegram.org"
const URL_GET_FILE = BASE_URL + "/bot" + BOT_TOKEN + "/getFile"
const URL_DOWNLOAD_FILE = BASE_URL + "/file/bot" + BOT_TOKEN
const GCS_BUCKET = PROJECT_ID + "-" + "telegram-object-sent"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "telegram_utils_log_file_downloader"

func init() {
	functions.HTTP("TelegramUtilsFileDownloader", TelegramUtilsFileDownloader)
}

func TelegramUtilsFileDownloader(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubSubscription PubsubSubscription
	var pubsubData PubsubData

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubSubscription); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode PubSub Message.
	if err := pubsubSubscription.decodePubSubData(&pubsubData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Download object from Telegram Server to GCS.
	if err := pubsubData.downloadFile(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
