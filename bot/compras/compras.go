package compras

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/faturas"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
)

// ProcessoAcoesCompras é responsável por coordenar as ações relacionadas a compras
func ProcessoAcoesCompras(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCartaoState *UserStateCompras) {

	switch message.Text {
	case "compras":
		gerarOpcoesAcoesCompras(bot, message)
	case "cadastrar_compra":
		cartoes := cartao.ListarCartoes(cartao.BaseURLCartoes)

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCartaoState)
	}
}

// gerarOpcoesAcoesCompras é responsável por gerar os botões para seleção das ações de compras para o usuário
func gerarOpcoesAcoesCompras(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	buttonOpcao1 := tgbotapi.NewKeyboardButton("cadastrar_compra")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("Opção 2")
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

// EnviarOpcoesCartoesFatura é responsável por enviar as opções de cartões disponíveis
func EnviarOpcoesCartoesFatura(bot *tgbotapi.BotAPI, chatID int64, cartao *cartao.ResPag, userStatesCompras *UserStateCompras) {
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

	step := "selecionar_fatura"

	userStatesCompras.CurrentStep = &step
}

// EnviarOpcoesFaturasCompras Função para enviar botões inline de seleção de faturas para o fluxo de cadastrar compra
func EnviarOpcoesFaturasCompras(bot *tgbotapi.BotAPI, chatID int64, faturas *faturas.ResPagFaturas, userStatesCompras *UserStateCompras, callbackQuery *tgbotapi.CallbackQuery) {
	res := cartao.BuscarCartao(cartao.BaseURLCartao, callbackQuery.Data)

	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Cartão Selecionado: %s", *res.Nome))
	edit.ReplyMarkup = nil

	_, err := bot.Send(edit)
	if err != nil {
		log.Panic(err)
	}

	// Criar slice para armazenar botões
	var buttons []tgbotapi.InlineKeyboardButton

	for _, invoice := range faturas.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*invoice.Nome, invoice.ID.String())

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

	step := "fatura_selecionada"

	userStatesCompras.CurrentStep = &step
}

// ListarComprasFatura é responsável por realizar a requisição de listagem para compras
func ListarComprasFatura(idFatura *string) (res ResComprasPag) {
	resp, err := http.Get(BaseURLCompras + "?fatura_id=" + *idFatura)
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

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

// ObterComprasPdf é responsável por realizar a requisição que obtém o pdf com as compras
func ObterComprasPdf(idFatura *uuid.UUID, idCartao *uuid.UUID) []byte {
	var (
		resp *http.Response
		err  error
	)

	if idFatura != nil && idCartao != nil {
		resp, err = http.Get(BaseURLComprasPdf + "?fatura_id=" + idFatura.String() + "&cartao_id=" + idCartao.String())
		if err != nil {
			fmt.Println("Erro ao fazer a requisição:", err)
			return nil
		}
		defer resp.Body.Close()
	} else {
		resp, err = http.Get(BaseURLComprasPdf + "?cartao_id=" + idCartao.String())
		if err != nil {
			fmt.Println("Erro ao fazer a requisição:", err)
			return nil
		}
		defer resp.Body.Close()
	}

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return nil
	}

	// Imprime a resposta da API
	fmt.Println("Resposta da API:", string(body))

	return body
}
