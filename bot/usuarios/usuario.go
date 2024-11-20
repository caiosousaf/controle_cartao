package usuarios

import (
	"bot_controle_cartao/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// GerarOpcoesAcoesUsuarios é responsável por enviar as opções de usuários disponíveis
func GerarOpcoesAcoesUsuarios(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	buttonOpcao1 := tgbotapi.NewKeyboardButton("Login")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("Cadastrar Usuario")

	keyboard := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{buttonOpcao1, buttonOpcao2},
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Selecione uma opção:")
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// Login é responsável por realizar o login do usuário
func Login(bot *tgbotapi.BotAPI, chatID int64, email *string, senha string, userTokens map[int64]string) error {
	var ambiente = utils.ValidarAmbiente()

	url := "/usuarios/login"
	requestBody := fmt.Sprintf(`{"email": "%s", "senha": "%s"}`, *email, senha)

	resp, err := http.Post(ambiente+url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	email = nil

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha no login: %s", resp.Status)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	token := response["token"].(string)

	userTokens[chatID] = token

	// Envie uma mensagem de sucesso ao usuário
	msg := tgbotapi.NewMessage(chatID, "Login realizado com sucesso!")
	bot.Send(msg)

	return nil
}
