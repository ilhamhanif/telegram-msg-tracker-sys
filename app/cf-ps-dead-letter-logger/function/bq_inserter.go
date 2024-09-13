package function

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

type BqRow struct {
	PubsubSubscription PubsubSubscription
}

func (r *BqRow) Save() (map[string]bigquery.Value, string, error) {

	/*
		A method to format BigQuery row record.
	*/

	tz, _ := time.LoadLocation("Asia/Jakarta")
	currDatetime := time.Now().In(tz)
	logDatetime := currDatetime.Format("2006-01-02T15:04:05")
	logEpoch := currDatetime.Format("20060102150405")
	logDate := currDatetime.Format("2006-01-02")

	ptRaw := r.PubsubSubscription.Message.Attributes.CloudPubSubDeadLetterSourceTopicPublishTime
	pt, _ := time.Parse(time.RFC3339, ptRaw)
	publishDatetime := pt.Format("2006-01-02T15:04:05")
	publishEpoch := pt.Format("20060102150405")
	publishDate := pt.Format("2006-01-02")

	dltRaw := r.PubsubSubscription.Message.PublishTime
	dlt, _ := time.Parse("2006-01-02T15:04:05.999999Z0700", dltRaw)
	deadLetterDatetime := dlt.Format("2006-01-02T15:04:05")
	deadLetterEpoch := dlt.Format("20060102150405")
	deadLetterDate := dlt.Format("2006-01-02")

	deliveryAttempt := r.PubsubSubscription.Message.Attributes.CloudPubSubDeadLetterSourceDeliveryCount
	subscriptionName := r.PubsubSubscription.Message.Attributes.CloudPubSubDeadLetterSourceSubscription
	subscriptionProjectID := r.PubsubSubscription.Message.Attributes.CloudPubSubDeadLetterSourceSubscriptionProject
	messageID := r.PubsubSubscription.Message.MessageID
	messageData := r.PubsubSubscription.Message.Data
	messageDataDecoded, _ := base64.StdEncoding.DecodeString(messageData)
	messageDataDecodedStr := string(messageDataDecoded)
	isRecycled := false

	return map[string]bigquery.Value{
		"delivery_attempt":        deliveryAttempt,
		"subscription_nm":         subscriptionName,
		"subscription_project_id": subscriptionProjectID,
		"publish_datetime":        publishDatetime,
		"publish_date":            publishDate,
		"publish_epoch":           publishEpoch,
		"message_id":              messageID,
		"message_data":            messageData,
		"message_data_decoded":    messageDataDecodedStr,
		"dead_letter_datetime":    deadLetterDatetime,
		"dead_letter_date":        deadLetterDate,
		"dead_letter_epoch":       deadLetterEpoch,
		"is_recycled":             isRecycled,
		"log_date":                logDate,
		"log_datetime":            logDatetime,
		"log_epoch":               logEpoch,
	}, bigquery.NoDedupeID, nil

}

func (r *BqRow) insertBqRows() error {

	/*
		A method to insert row records to GCP BigQuery.
	*/

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return fmt.Errorf("insertBqRows: Error initiating BigQuery client: %w", err)
	}
	defer client.Close()

	inserter := client.Dataset(BQ_DATASET_NAME).Table(BQ_TABLE_NAME).Inserter()
	if err := inserter.Put(ctx, []*BqRow{r}); err != nil {
		return fmt.Errorf("insertBqRows: Error inserting rows to BigQuery: %w", err)
	}

	return nil

}
