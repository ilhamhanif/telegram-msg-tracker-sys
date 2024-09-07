package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
)

type BqRow struct {
	IdentificationResult IdentificationResult
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

	updateId := r.IdentificationResult.UpdateId
	updateEpoch := r.IdentificationResult.UpdateEpoch
	updateDate := r.IdentificationResult.UpdateDate
	updateDatetime := r.IdentificationResult.UpdateDatetime
	result, _ := json.Marshal(r.IdentificationResult.Result)
	resultStr := string(result)

	return map[string]bigquery.Value{
		"update_id":       updateId,
		"update_epoch":    updateEpoch,
		"update_date":     updateDate,
		"update_datetime": updateDatetime,
		"result":          resultStr,
		"log_date":        logDate,
		"log_datetime":    logDatetime,
		"log_epoch":       logEpoch,
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
