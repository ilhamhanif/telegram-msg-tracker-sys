package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const PROJECT_ID = "protean-quanta-434205-p5"
const PUBSUB_TOPIC_IDENTIFICATOR = "tlgrm_msg_identificator"
const PUBSUB_TOPIC_LOGGER = "tlgrm_msg_upd_logger"

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
	if err := telegramMsgUpdate.publishToPubSub(PUBSUB_TOPIC_IDENTIFICATOR); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the `raw` message to LOGGER through Pub/Sub.
	if err := telegramMsgUpdate.publishToPubSub(PUBSUB_TOPIC_LOGGER); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
