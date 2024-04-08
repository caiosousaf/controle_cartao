package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// EnviaMensagem é responsável por encapsular a função Send da libi
func EnviaMensagem(bot *tgbotapi.BotAPI, c tgbotapi.Chattable) {
	_, err := bot.Send(c)
	if err != nil {
		log.Panic(err)
	}
}
