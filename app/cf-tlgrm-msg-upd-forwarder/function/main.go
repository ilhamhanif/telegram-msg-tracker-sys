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
	functions.HTTP("TelegramMsgUpdateForwarder", TelegramMsgUpdateForwarder)
}

func TelegramMsgUpdateForwarder(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var telegramMsgUpdate TelegramApiModelUpdate

	// Receive and parse HTTP push data message
	// from Telegram Webhook.
	if err := json.NewDecoder(r.Body).Decode(&telegramMsgUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the message to IDENTIFICATOR through Pub/Sub.
	if err := telegramMsgUpdate.publishToPubSub(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
