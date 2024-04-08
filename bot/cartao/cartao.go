package cartao

import (
	"bot_controle_cartao/compras"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ProcessoAcoesCartoes é responsável por coordenar as ações relacionadas a cartões
func ProcessoAcoesCartoes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCartaoState *UserStateCartao) {

	switch message.Text {
	case "cartoes":
		gerarOpcoesAcoesCartao(bot, message)
	case "extrato":
		cartoes := ListarCartoes(BaseURLCartoes)

		gerarOpcoesCartoesDisponiveis(bot, message.Chat.ID, &cartoes, userCartaoState)
		//case "cadastrar_cartao":
		//	inicioCriacaoFatura(bot, message.Chat.ID, userStates)
		//	userCartaoState.CurrentStepBool = true
	}
}

// ProcessarCasosStepExtratoCartao é a função responsável por controlar os steps do usuário para o fluxo de extrato de um cartão
func ProcessarCasosStepExtratoCartao(userStatesCartao *UserStateCartao, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch userStatesCartao.CurrentStep {
	case "selecionar_ano":
		userStatesCartao.NovoCartaoData.ID = update.CallbackQuery.Data

		EnviarOpcoesAno(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userStatesCartao)
	case "ano_selecionado":
		idCartaoUUID, err := uuid.Parse(userStatesCartao.NovoCartaoData.ID)
		if err != nil {
			log.Panic(err)
		}

		edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Cartão Selecionado: %s", update.CallbackQuery.Data))
		edit.ReplyMarkup = nil

		_, err = bot.Send(edit)
		if err != nil {
			log.Panic(err)
		}

		pdfContent := compras.ObterComprasPdf(nil, &idCartaoUUID)

		msgPdfCompras := tgbotapi.NewDocumentUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileReader{
			Name:   "compras_" + update.CallbackQuery.Data + ".pdf",
			Reader: bytes.NewBuffer(pdfContent),
			Size:   int64(len(pdfContent)),
		})

		_, err = bot.Send(msgPdfCompras)
		if err != nil {
			log.Panic(err)
		}
	}
}

// gerarOpcoesAcoesCartao é responsável por gerar os botões para seleção das ações de cartões para o usuário
func gerarOpcoesAcoesCartao(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	buttonOpcao1 := tgbotapi.NewKeyboardButton("extrato")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("cadastrar_cartao")
	buttonOpcao3 := tgbotapi.NewKeyboardButton("Opção 3")
	buttonOpcao4 := tgbotapi.NewKeyboardButton("Opção 4")

	keyboard := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{buttonOpcao1, buttonOpcao2},
		[]tgbotapi.KeyboardButton{buttonOpcao3, buttonOpcao4},
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, "Selecione uma opção:")
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// gerarOpcoesCartoesDisponiveis Função para enviar botões inline de seleção de cartões
func gerarOpcoesCartoesDisponiveis(bot *tgbotapi.BotAPI, chatID int64, cartao *ResPag, userCartaoState *UserStateCartao) {
	// Criar slice para armazenar botões
	var buttons []tgbotapi.InlineKeyboardButton

	for _, card := range cartao.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*card.Nome, card.ID.String())

		buttons = append(buttons, button)
	}

	lenOptions1, lenOptions2 := len(buttons)/2, len(buttons)/2
	if len(buttons)%2 != 0 {
		lenOptions1++ // Adiciona 1 à primeira parte se o número de botões for ímpar
	}

	buttonsOne := make([]tgbotapi.InlineKeyboardButton, lenOptions1)
	buttonsTwo := make([]tgbotapi.InlineKeyboardButton, lenOptions2)

	copy(buttonsOne, buttons[:lenOptions1])
	copy(buttonsTwo, buttons[lenOptions2:])

	// Criar teclado inline com os botões
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		buttonsOne,
		buttonsTwo,
	)

	// Configurar a mensagem com o teclado
	msg := tgbotapi.NewMessage(chatID, "Selecione um cartão: ")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "selecionar_ano"

	userCartaoState.CurrentStep = step
}

// EnviarOpcoesAno envia as opções para que o usuário selecione o ano que será usado quando for gerado o extrato
func EnviarOpcoesAno(bot *tgbotapi.BotAPI, chatID int64, callbackQuery *tgbotapi.CallbackQuery, userCartaoState *UserStateCartao) {
	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Opção Selecionada: %s", callbackQuery.Data))
	edit.ReplyMarkup = nil

	_, err := bot.Send(edit)
	if err != nil {
		log.Panic(err)
	}

	// Criar slice para armazenar botões
	var buttons []tgbotapi.InlineKeyboardButton

	currentYear := time.Now().Year()

	// Array para armazenar os últimos cinco anos
	var lastFiveYears [5]int

	// Itera de trás para frente, cinco anos
	for i := 0; i < 5; i++ {
		// Calcula o ano retroativo
		year := currentYear - i
		// Armazena o ano no array
		lastFiveYears[i] = year
	}

	// Imprime o array
	fmt.Println(lastFiveYears)

	for _, anoPossivel := range lastFiveYears {
		anoPossivelString := strconv.Itoa(anoPossivel)
		button := tgbotapi.NewInlineKeyboardButtonData(anoPossivelString, anoPossivelString)

		buttons = append(buttons, button)
	}

	lenOptions1, lenOptions2 := len(buttons)/2, len(buttons)/2
	if len(buttons)%2 != 0 {
		lenOptions1++
		lenOptions2++
	}

	buttonsOne := make([]tgbotapi.InlineKeyboardButton, lenOptions1)
	buttonsTwo := make([]tgbotapi.InlineKeyboardButton, lenOptions2-1)

	copy(buttonsOne, buttons[:lenOptions1])
	copy(buttonsTwo, buttons[lenOptions2:])

	// Criar teclado inline com os botões
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		buttonsOne,
		buttonsTwo,
	)

	// Configurar a mensagem com o teclado
	msg := tgbotapi.NewMessage(chatID, "Selecione uma fatura:")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "ano_selecionado"

	userCartaoState.CurrentStep = step
}

// ListarCartoes é responsável por listar os cartões cadastrados
func ListarCartoes(url string) (cartoes ResPag) {
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

	if err := json.Unmarshal(body, &cartoes); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}
