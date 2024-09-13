package function

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

type DownloadedFile struct {
	FileUpdateID   int    `json:"file_update_id"`
	FileUpdateDate string `json:"file_update_date"`
	FileName       string `json:"file_name"`
	FilePath       string `json:"file_path"`
	FilePathRel    string `json:"file_path_rel"`
}

func (d *DownloadedFile) uploadFileToGCS() error {

	/*
		A function to orchestrate (create and delete) an object in Google GCS
	*/

	// Setting up GCS client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("uploadFileToGCS: Error when setting up GCS client: %w", err)
	}
	defer client.Close()

	// Get local object
	f, err := os.Open(d.FilePath)
	if err != nil {
		return fmt.Errorf("uploadFileToGCS: Error opening local file: %w", err)
	}
	defer f.Close()

	// Setup automatic cancel after 50 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	o := client.Bucket(GCS_BUCKET).Object(
		filepath.Join(
			"update_date="+d.FileUpdateDate,
			"update_id="+strconv.Itoa(d.FileUpdateID),
			d.FilePathRel,
		),
	)
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("uploadFileToGCS: Error uploading local file to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("uploadFileToGCS: Error closing local file: %w", err)
	}

	return nil

}
