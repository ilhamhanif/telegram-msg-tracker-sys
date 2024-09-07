package function

import (
	"fmt"
	"strings"

	"github.com/go-telegram/bot/models"
)

type IdentificationResult struct {
	UpdateId       int
	UpdateEpoch    int
	UpdateDate     string
	UpdateDatetime string
	Result         struct {
		Type   string
		ChatId int64
		Text   string
		Photo  []models.PhotoSize
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

	if r.Result.Type == "BOT COMMAND" {
		if !strings.HasPrefix(r.Result.Text, "/") {
			botSendMessageParams.ChatID = r.Result.ChatId
			botSendMessageParams.Text = "Bot command have to be started with /."
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
