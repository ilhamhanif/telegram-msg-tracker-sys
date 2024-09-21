package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

type BqRow struct {
	PubsubSubscription PubsubSubscription `json:"pubsub_subscription"`
	PubsubData         PubsubData         `json:"pubsub_data"`
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

	edtRaw := r.PubsubSubscription.Message.Attributes.EventTime
	edt, _ := time.Parse("2006-01-02T15:04:05.999999Z0700", edtRaw)
	eventDatetime := edt.Format("2006-01-02T15:04:05")
	eventEpoch := edt.Format("20060102150405")
	eventDate := edt.Format("2006-01-02")

	udtRaw := r.PubsubData.Updated
	udt, _ := time.Parse("2006-01-02T15:04:05.999999Z0700", udtRaw)
	updatedDatetime := udt.Format("2006-01-02T15:04:05")
	updatedEpoch := udt.Format("20060102150405")
	updatedDate := udt.Format("2006-01-02")

	cdtRaw := r.PubsubData.TimeCreated
	cdt, _ := time.Parse("2006-01-02T15:04:05.999999Z0700", cdtRaw)
	createdDatetime := cdt.Format("2006-01-02T15:04:05")
	createdEpoch := cdt.Format("20060102150405")
	createdDate := cdt.Format("2006-01-02")

	sudtRaw := r.PubsubData.TimeStorageClassUpdated
	sudt, _ := time.Parse("2006-01-02T15:04:05.999999Z0700", sudtRaw)
	storageUpdatedDatetime := sudt.Format("2006-01-02T15:04:05")
	storageUpdatedEpoch := sudt.Format("20060102150405")
	storageUpdatedDate := sudt.Format("2006-01-02")

	raw, _ := json.Marshal(r)
	rawStr := string(raw)

	return map[string]bigquery.Value{
		"bucket_id":                      r.PubsubSubscription.Message.Attributes.BucketID,
		"event_datetime":                 eventDatetime,
		"event_date":                     eventDate,
		"event_epoch":                    eventEpoch,
		"event_type":                     r.PubsubSubscription.Message.Attributes.EventType,
		"notification_config":            r.PubsubSubscription.Message.Attributes.NotificationConfig,
		"object_generation":              r.PubsubSubscription.Message.Attributes.ObjectGeneration,
		"payload_format":                 r.PubsubSubscription.Message.Attributes.PayloadFormat,
		"object_id":                      r.PubsubSubscription.Message.Attributes.ObjectID,
		"kind":                           r.PubsubData.Kind,
		"id":                             r.PubsubData.ID,
		"self_link":                      r.PubsubData.SelfLink,
		"name":                           r.PubsubData.Name,
		"metageneration":                 r.PubsubData.MetaGeneration,
		"content_type":                   r.PubsubData.ContentType,
		"created_datetime":               createdDatetime,
		"created_date":                   createdDate,
		"created_epoch":                  createdEpoch,
		"updated_datetime":               updatedDatetime,
		"updated_date":                   updatedDate,
		"updated_epoch":                  updatedEpoch,
		"storage_class":                  r.PubsubData.StorageClass,
		"storage_class_updated_datetime": storageUpdatedDatetime,
		"storage_class_updated_date":     storageUpdatedDate,
		"storage_class_updated_epoch":    storageUpdatedEpoch,
		"size":                           r.PubsubData.Size,
		"media_link":                     r.PubsubData.MediaLink,
		"log_date":                       logDate,
		"log_datetime":                   logDatetime,
		"log_epoch":                      logEpoch,
		"raw":                            rawStr,
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
