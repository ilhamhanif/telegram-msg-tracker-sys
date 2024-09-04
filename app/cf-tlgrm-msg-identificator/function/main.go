package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

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

	// Return `ok`
	fmt.Fprint(w, "ok")

}
