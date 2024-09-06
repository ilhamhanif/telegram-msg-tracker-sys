package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-telegram/bot/models"
)

type PubsubData models.Update

func init() {
	functions.HTTP("TelegramMsgIdentificator", TelegramMsgIdentificator)
}

func TelegramMsgIdentificator(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubMessage PubsubSubscription
	var pubsubData PubsubData
	var identificationResult IdentificationResult

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubMessage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode PubSub data to get the `raw` data.
	if err := pubsubMessage.decodePubSubData(&pubsubData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the `raw` message to LOGGER through Pub/Sub.
	if err := pubsubData.publishToPubSub(); err != nil {
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

	// Return `ok`.
	fmt.Fprint(w, "ok")

}
