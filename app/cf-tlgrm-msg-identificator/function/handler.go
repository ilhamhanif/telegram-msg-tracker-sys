package function

import (
	"strings"

	"github.com/go-telegram/bot"
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

type MessageSendMessage bot.SendMessageParams

func (r *IdentificationResult) check() error {

	/*
		A method to check all available scenario.
	*/

	if err := r.botCommand(); err != nil {
		return err
	} else if err := r.text(); err != nil {
		return err
	} else if err := r.photo(); err != nil {
		return err
	}

	return nil

}

func (r *IdentificationResult) botCommand() error {

	/*
		A method to handle Bot Command
		If Bot command is not started with /: send Error Message to source.
	*/

	if r.Result.Type == "BOT COMMAND" {
		if !strings.HasPrefix(r.Result.Text, "/") {
			return nil
		} else {
			return nil
		}
	}

	return nil
}

func (r *IdentificationResult) text() error {

	/*
		A method to handle Bot Command.
	*/

	return nil
}

func (r *IdentificationResult) photo() error {

	/*
		A method to handle Bot Command.
	*/

	return nil
}