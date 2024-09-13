package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-telegram/bot/models"
)

type PubsubData struct {
	UpdateID       int                 `json:"update_id"`
	UpdateEpoch    int                 `json:"update_epoch"`
	UpdateDate     string              `json:"update_date"`
	UpdateDatetime string              `json:"update_datetime"`
	Files          []map[string]string `json:"files"`
}

type ApiResult struct {
	Ok     bool        `json:"ok"`
	Result models.File `json:"result"`
}

func (pd *PubsubData) downloadFile() error {

	/*
		A method to download a File.
	*/

	// Initialize currDir and targetDir
	currDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("downloadFile: Failed to get current workdir: %w", err)
	}

	targetDir := filepath.Join(currDir, "target")
	if _, err = os.Stat(targetDir); os.IsNotExist(err) {
		err = os.Mkdir(targetDir, 0777)
		if err != nil {
			return fmt.Errorf("downloadFile: Failed to create target directory: %w", err)
		}
	}
	defer os.RemoveAll(targetDir)

	for _, v := range pd.Files {

		var apiResult ApiResult
		var downloadedFile DownloadedFile

		// Sent the data to Telegram API endpoint
		// using HTTP GET.
		gf_resp, gf_err := http.Get(URL_GET_FILE + "?file_id=" + v["file_id"])
		if gf_err != nil {
			return fmt.Errorf("downloadFile: Failed to send message through GET_FILE URL: %w", gf_err)
		}
		defer gf_resp.Body.Close()

		// Store the response status code and message
		body, err := io.ReadAll(gf_resp.Body)
		if err != nil {
			return fmt.Errorf("downloadFile: Failed to read API body response: %w", err)
		}
		if err := json.Unmarshal(body, &apiResult); err != nil {
			return fmt.Errorf("sendMessage: Failed to store API result: %w", err)
		}

		// Download file from Telegram
		// using HTTP GET.
		dl_resp, dl_err := http.Get(URL_DOWNLOAD_FILE + "/" + apiResult.Result.FilePath)
		if dl_err != nil {
			return fmt.Errorf("downloadFile: Failed to send message through DOWNLOAD_FILE URL: %w", err)
		}
		defer dl_resp.Body.Close()

		// Save file in designated path.
		fileName := apiResult.Result.FileID + "_" + apiResult.Result.FileUniqueID + "_" + strings.Replace(apiResult.Result.FilePath, "/", "", -1)
		filePath := filepath.Join(targetDir, fileName)
		filePathRel, _ := filepath.Rel(targetDir, filePath)
		out, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("downloadFile: Failed to create target file: %w", err)
		}
		defer out.Close()

		if _, err = io.Copy(out, dl_resp.Body); err != nil {
			return fmt.Errorf("downloadFile: Failed to copy downloaded file into target directory: %w", err)
		}

		// Upload file to GCS
		downloadedFile.FileUpdateID = pd.UpdateID
		downloadedFile.FileUpdateDate = pd.UpdateDate
		downloadedFile.FileName = fileName
		downloadedFile.FilePath = filePath
		downloadedFile.FilePathRel = filePathRel
		if err = downloadedFile.uploadFileToGCS(); err != nil {
			return fmt.Errorf("downloadFile: Failed to upload file to GCS: %w", err)
		}

		// Insert log to Google BigQuery.
		var bqRows = BqRow{
			PubsubData:     *pd,
			DownloadedFile: downloadedFile,
		}
		if err := bqRows.insertBqRows(); err != nil {
			return fmt.Errorf("downloadFile: Failed to save log to BigQuery: %w", err)
		}

		// Remove file after successfully upload to GCS
		if err = os.Remove(filePath); err != nil {
			return fmt.Errorf("downloadedFile: Failed to remove downloaded file: %w", err)
		}
	}

	return nil

}
