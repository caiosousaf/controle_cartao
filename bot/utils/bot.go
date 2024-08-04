package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

// CancelarOperacao é uma função que cancela uma operação com steps e redefine o step atual
func CancelarOperacao(bot *tgbotapi.BotAPI, mensagem, currentStep *string, chatID int64) bool {
	if strings.ToLower(*mensagem) == "cancelar" {
		currentStep = nil
		msg := tgbotapi.NewMessage(chatID, "Operação cancelada.")
		EnviaMensagem(bot, msg)
		return true
	}

	return false
}
