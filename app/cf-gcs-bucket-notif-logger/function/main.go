package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const PROJECT_ID = "protean-quanta-434205-p5"
const BQ_DATASET_NAME = "ops"
const BQ_TABLE_NAME = "gcs_bucket_notif_log"

func init() {
	functions.HTTP("GCSBucketNotifLogger", GCSBucketNotifLogger)
}

func GCSBucketNotifLogger(w http.ResponseWriter, r *http.Request) {

	/*
		Main Function.
	*/

	var pubsubSubscription PubsubSubscription
	var pubsubData PubsubData

	// Receive and parse GCP Pub/Sub HTTP push data message.
	if err := json.NewDecoder(r.Body).Decode(&pubsubSubscription); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	// Decode PubSub data to get the `raw` data.
	if err := pubsubSubscription.decodePubSubData(&pubsubData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert data from message received to Google BigQuery
	var bqRows = BqRow{
		PubsubSubscription: pubsubSubscription,
		PubsubData:         pubsubData,
	}
	if err := bqRows.insertBqRows(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Always return success
	fmt.Fprint(w, "ok")
}
