package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (m *PubsubData) sendMessage(v *ApiResult) error {

	/*
		A method to send a message with Telegram API.
	*/

	// Setup message in JSON.
	messageJson, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("sendMessage: Error: %w", err)
	}

	// Sent the data to Telegram API endpoint
	// using HTTP POST.
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(messageJson))
	if err != nil {
		return fmt.Errorf("sendMessage: Failed to setup a new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sendMessage: Failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Store the response status code and message
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("sendMessage: Failed to read API body response: %w", err)
	}
	if err := json.Unmarshal(body, &v.Message); err != nil {
		return fmt.Errorf("sendMessage: Failed to store API result: %w", err)
	}
	v.StatusCode = resp.StatusCode

	return nil
}
