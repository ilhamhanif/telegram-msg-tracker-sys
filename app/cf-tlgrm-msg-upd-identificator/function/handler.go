package function

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
)

type IdentificationResult struct {
	UpdateID       int    `json:"update_id"`
	UpdateEpoch    int    `json:"update_epoch"`
	UpdateDate     string `json:"update_date"`
	UpdateDatetime string `json:"update_datetime"`
	Result         struct {
		Type   string             `json:"type"`
		ChatID int64              `json:"chat_id"`
		Text   string             `json:"text"`
		Photo  []models.PhotoSize `json:"photos"`
	}
}

func (r *IdentificationResult) check() error {

	/*
		A method to check all available scenario.
	*/

	if err := r.isBotCommand(); err != nil {
		return fmt.Errorf("check: %w", err)
	} else if err := r.isText(); err != nil {
		return fmt.Errorf("check: %w", err)
	} else if err := r.isPhoto(); err != nil {
		return fmt.Errorf("check: %w", err)
	}

	return nil

}

func (r *IdentificationResult) isBotCommand() error {

	/*
		A method to handle Bot Command
		1. If Bot command is not started with /: send Error Message to source.
	*/

	var botSendMessageParams BotSendMessageParams

	if r.Result.Type == "BOT_COMMAND" {
		if !strings.HasPrefix(r.Result.Text, "/") {
			botSendMessageParams.UpdateID = r.UpdateID
			botSendMessageParams.UpdateEpoch = r.UpdateEpoch
			botSendMessageParams.UpdateDate = r.UpdateDate
			botSendMessageParams.UpdateDatetime = r.UpdateDatetime
			botSendMessageParams.Params.ChatID = r.Result.ChatID
			botSendMessageParams.Params.Text = "Bot command have to be started with /."
		}
	}
	if err := botSendMessageParams.sendMessage(); err != nil {
		return fmt.Errorf("isBotCommand: Error sending to PubSub: %w", err)
	}

	return nil

}

func (r *IdentificationResult) isText() error {

	/*
		A method to handle Text.
	*/

	return nil

}

func (r *IdentificationResult) isPhoto() error {

	/*
		A method to handle Photo.
	*/

	return nil

}
