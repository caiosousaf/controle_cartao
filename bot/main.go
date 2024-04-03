package main

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/faturas"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)

func realizarGetString(url string) (msgGet string) {
	// Realiza a requisição GET para a API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	// Imprime a resposta da API
	fmt.Println("Resposta da API:", string(body))

	var cartoes cartao.ResPag
	if err := json.Unmarshal(body, &cartoes); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	for _, cartao := range cartoes.Dados {
		msgGet += fmt.Sprintf("<i>ID: %s\n\n</i>", cartao.ID.String())
		msgGet += fmt.Sprintf("<i>Nome: %s\n\n</i>", *cartao.Nome)
		msgGet += fmt.Sprintf("<i>Data de Criação: %s\n\n</i>", cartao.DataCriacao.String())
		if cartao.DataDesativacao != nil {
			msgGet += fmt.Sprintf("<i>Data de Desativação: %s\n\n</i>", cartao.DataDesativacao.String())
		}
		msgGet += "------\n\n"
	}

	return
}

func main() {

	// Inicialize o token do seu bot aqui
	token := "6821239738:AAGyxhdn27UYG7TSm31DpS_cKo0ezbzoySA"

	// Cria um novo bot com o token fornecido
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Configuração de debug para receber informações detalhadas
	bot.Debug = true

	log.Printf("Autorizado como %s", bot.Self.UserName)

	userStates := make(map[int64]*faturas.UserState)
	userCompraFaturas := &faturas.UserStepComprasFatura{
		Cartoes: []string{}, // Preencha a fatia de cartões conforme necessário
		Opcao:   nil,        // Ou atribua um valor ao ponteiro de opção, se necessário
	}

	// Configuração de atualização com o webhook ou polling
	// Aqui, estamos usando a opção de polling para obter atualizações
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	// Loop pelas atualizações recebidas do bot
	for update := range updates {
		if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Message.Text)

			switch *userCompraFaturas.Opcao {
			case "fatura_selecionada":
				faturasCartao := faturas.ListarFaturas(fmt.Sprintf(faturas.BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data))

				faturas.EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, update.CallbackQuery)
			case "cartao_fatura_selecionado":
				faturas.ProcessCallbackQuery(bot, update.CallbackQuery)
			}
		} else if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Se a mensagem do usuário for "/start", envie uma mensagem de boas-vindas
			if update.Message.Text == "/start" {
				// Criando um teclado de resposta
				buttonOpcao1 := tgbotapi.NewKeyboardButton("cartoes")
				buttonOpcao2 := tgbotapi.NewKeyboardButton("faturas")
				buttonOpcao3 := tgbotapi.NewKeyboardButton("Opção 3")
				buttonOpcao4 := tgbotapi.NewKeyboardButton("Opção 4")

				keyboard := tgbotapi.NewReplyKeyboard(
					[]tgbotapi.KeyboardButton{buttonOpcao1, buttonOpcao2},
					[]tgbotapi.KeyboardButton{buttonOpcao3, buttonOpcao4},
				)

				// Configurando a mensagem de boas-vindas com o teclado de resposta
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Selecione uma opção:")
				msg.ReplyMarkup = keyboard

				// Enviando a mensagem
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
				faturas.AcaoAnterior = "start"
			}

			if update.Message.Text == "/cartoes" || update.Message.Text == "cartoes" {
				msgGet := realizarGetString(faturas.BaseURLCartoes)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgGet)
				msg.ParseMode = "HTML"
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
				faturas.AcaoAnterior = "cartoes"
			}

			if update.Message.Text == "faturas" || update.Message.Text == "/faturas" || faturas.AcaoAnterior == "faturas" {
				faturas.ProcessoAcoesFaturas(bot, update.Message, userStates, userCompraFaturas)

				faturas.AcaoAnterior = "faturas"
			}
		}
	}
}
