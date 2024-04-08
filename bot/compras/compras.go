package compras

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/categorias"
	"bot_controle_cartao/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// ProcessoAcoesCadastroCompra é responsável pelo fluxo de cadastro de uma compra
func ProcessoAcoesCadastroCompra(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCompras *UserStateCompras) {
	if strings.ToLower(message.Text) == "cancelar" {
		userCompras.CurrentStep = nil
		userCompras.NovaCompraData = NovaCompra{
			Nome:               nil,
			Descricao:          nil,
			LocalCompra:        nil,
			CategoriaID:        nil,
			ValorParcela:       nil,
			ParcelaAtual:       nil,
			QuantidadeParcelas: nil,
			DataCompra:         nil,
		}

		return
	}

	switch *userCompras.CurrentStep {
	case "inicio_cadastro_compra":
		userCompras.NovaCompraData.Nome = &message.Text
		*userCompras.CurrentStep = "nome_selecionado"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Digite a descrição da compra:")

		utils.EnviaMensagem(bot, msg)
	case "nome_selecionado":
		userCompras.NovaCompraData.Descricao = &message.Text
		*userCompras.CurrentStep = "descricao_selecionada"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga qual foi o local de compra:")

		utils.EnviaMensagem(bot, msg)
	case "descricao_selecionada":
		userCompras.NovaCompraData.LocalCompra = &message.Text
		*userCompras.CurrentStep = "local_selecionado"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga qual foi o valor da parcela/compra:")

		utils.EnviaMensagem(bot, msg)
	case "local_selecionado":
		valorParcela, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			log.Panic(err)
		}

		userCompras.NovaCompraData.ValorParcela = &valorParcela
		*userCompras.CurrentStep = "valor_selecionado"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga qual a parcela atual da compra, ex: 1, 2")

		utils.EnviaMensagem(bot, msg)
	case "valor_selecionado":
		parcelaAtual, err := strconv.ParseInt(message.Text, 10, 64)
		if err != nil {
			log.Panic(err)
		}

		userCompras.NovaCompraData.ParcelaAtual = &parcelaAtual
		*userCompras.CurrentStep = "parcela_selecionada"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga a quantidade de parcelas:")

		utils.EnviaMensagem(bot, msg)
	case "parcela_selecionada":
		qtdParcelas, err := strconv.ParseInt(message.Text, 10, 64)
		if err != nil {
			log.Panic(err)
		}

		userCompras.NovaCompraData.QuantidadeParcelas = &qtdParcelas
		*userCompras.CurrentStep = "quantidade_parcelas_selecionada"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga qual foi a data da compra, ex: 2024-02-20")

		utils.EnviaMensagem(bot, msg)
	case "quantidade_parcelas_selecionada":
		userCompras.NovaCompraData.DataCompra = &message.Text
		*userCompras.CurrentStep = ""

		CadastrarCompra(userCompras)
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

// EnviarOpcoesCategoriasCompras é responsável por envia via telegram as categorias das compras
func EnviarOpcoesCategoriasCompras(bot *tgbotapi.BotAPI, chatID int64, categorias *categorias.ResCategoriasPag, userStatesCompras *UserStateCompras, callbackQuery *tgbotapi.CallbackQuery) {
	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Fatura selecionada!"))
	edit.ReplyMarkup = nil

	_, err := bot.Send(edit)
	if err != nil {
		log.Panic(err)
	}

	userStatesCompras.FaturaID = &callbackQuery.Data
	// Criar slice para armazenar botões
	var buttons []tgbotapi.InlineKeyboardButton

	for _, categoria := range categorias.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*categoria.Nome, categoria.ID.String())

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
	msg := tgbotapi.NewMessage(chatID, "Selecione uma categoria: ")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "categoria_selecionada"

	userStatesCompras.CurrentStep = &step
}

// InicioCriacaoCompra é responsável por iniciar o processo de criação de uma compra
func InicioCriacaoCompra(bot *tgbotapi.BotAPI, chatID int64, callbackQuery *tgbotapi.CallbackQuery, userCompras *UserStateCompras) {
	msg := tgbotapi.NewMessage(chatID, "Em qualquer momento do cadastro, digite 'cancelar' para cancelar o cadastro")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	categoriaId, err := uuid.Parse(callbackQuery.Data)
	if err != nil {
		log.Panic(err)
	}

	userCompras.NovaCompraData.CategoriaID = &categoriaId
	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Categoria selecionada!"))
	edit.ReplyMarkup = nil

	_, err = bot.Send(edit)
	if err != nil {
		log.Panic(err)
	}

	*userCompras.CurrentStep = "inicio_cadastro_compra"

	msg = tgbotapi.NewMessage(chatID, "Por favor, insira o nome da Compra:")
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
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

// CadastrarCompra é responsável por realizar a requisição para cadastrar uma compra
func CadastrarCompra(compras *UserStateCompras) (res ResComprasPag) {
	const baseURLCadastroCompras = "http://localhost:8080/cadastros/fatura"
	// Montar os dados a serem enviados no corpo do POST
	dados := NovaCompra{
		Nome:               compras.NovaCompraData.Nome,
		Descricao:          compras.NovaCompraData.Descricao,
		LocalCompra:        compras.NovaCompraData.LocalCompra,
		CategoriaID:        compras.NovaCompraData.CategoriaID,
		ValorParcela:       compras.NovaCompraData.ValorParcela,
		ParcelaAtual:       compras.NovaCompraData.ParcelaAtual,
		QuantidadeParcelas: compras.NovaCompraData.QuantidadeParcelas,
		DataCompra:         compras.NovaCompraData.DataCompra,
	}

	// Codificar os dados em formato JSON
	dadosJSON, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao codificar os dados JSON:", err)
		return
	}

	// Fazer a requisição POST
	resp, err := http.Post(baseURLCadastroCompras+fmt.Sprintf("/%s/compras", *compras.FaturaID), "application/json", bytes.NewBuffer(dadosJSON))
	if err != nil {
		fmt.Println("Erro ao fazer a requisição POST:", err)
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
