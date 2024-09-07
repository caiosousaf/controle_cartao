package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
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

// ValidarAmbiente é uma função que valida se a aplicação está em produção
func ValidarAmbiente() string {
	sistemaProducao := os.Getenv("PROD")

	if sistemaProducao == "" {
		sistemaLocal := os.Getenv("LOCAL")

		return sistemaLocal
	}

	return sistemaProducao
}
