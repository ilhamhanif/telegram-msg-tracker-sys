package function

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/bigquery"
)

type BqRow struct {
	UpdateMessage TelegramApiModelUpdate
}

func (r *BqRow) Save() (map[string]bigquery.Value, string, error) {

	/*
		A method to format BigQuery row record
	*/

	tz, _ := time.LoadLocation("Asia/Jakarta")
	currDatetime := time.Now().In(tz)
	logDatetime := currDatetime.Format("2006-01-02T15:04:05")
	logEpoch := currDatetime.Format("20060102150405")
	logDate := currDatetime.Format("2006-01-02")
	updateId := r.UpdateMessage.ID
	update, _ := json.Marshal(r.UpdateMessage)
	updateStr := string(update)

	return map[string]bigquery.Value{
		"update_id":    updateId,
		"update":       updateStr,
		"log_date":     logDate,
		"log_datetime": logDatetime,
		"log_epoch":    logEpoch,
	}, bigquery.NoDedupeID, nil

}

func insertBqRows(rows []*BqRow) error {

	/*
		A function to insert row records to GCP BigQuery
	*/

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, PROJECT_ID)
	if err != nil {
		return err
	}
	defer client.Close()

	inserter := client.Dataset(BQ_DATASET_NAME).Table(BQ_TABLE_NAME).Inserter()
	if err := inserter.Put(ctx, rows); err != nil {
		return err
	}

	return nil

}
