package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

type TelegramApiModelUpdate models.Update

func init() {
	functions.HTTP("TelegramMsgUpdateLogger", TelegramMsgUpdateLogger)
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

	// Decode PubSub Message
	if err := pubsubMessage.decodePubSubData(&telegramMsgUpdate); err != nil {
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
