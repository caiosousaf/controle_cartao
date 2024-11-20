package faturas

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/compras"
	"bot_controle_cartao/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// ProcessoAcoesFaturas é responsável por ter todos os processos das faturas
func ProcessoAcoesFaturas(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userStates map[int64]*UserState, userCompraFatura *UserStepComprasFatura, userTokens map[int64]string) {
	if userState, ok := userStates[message.Chat.ID]; ok {
		if userState.CurrentStepBool {
			finalizarCadastroFatura(bot, message, userTokens, userState)
		}
	}

	switch message.Text {
	case "Faturas":
		gerarOpcoesFatura(bot, message)
	case "Compras Fatura":
		cartoes, err := cartao.ListarCartoes(cartao.BaseURLCartoes, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCompraFatura)
	case "Cadastrar Fatura":
		inicioCriacaoFatura(bot, message, userStates, userTokens, userCompraFatura)
		userStates[message.Chat.ID].CurrentStepBool = true
	case "Pagar Fatura":
		cartoes, err := cartao.ListarCartoes(cartao.BaseURLCartoes, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCompraFatura)

		*userCompraFatura.Opcao = "Pagar Fatura"
	}
}

// EnviaMensagemBoasVindas é responsável por enviar uma mensagem de boas vindas
func EnviaMensagemBoasVindas(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Bem-vindo ao Bot de Faturas! A qualquer momento digite 'cancelar' para cancelar operações de cadastro")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// inicioCriacaoFatura é responsável por iniciar o processo de criação de uma fatura
func inicioCriacaoFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userStates map[int64]*UserState, userTokens map[int64]string, userCompraFatura *UserStepComprasFatura) {
	userStates[message.Chat.ID] = &UserState{
		ChatID:      message.Chat.ID,
		CurrentStep: "cadastro_fatura",
	}

	cartoes, err := cartao.ListarCartoes(cartao.BaseURLCartoes, userTokens, message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCompraFatura)
}

// ContinuaCriacaoFatura é responsável por continuar o processo de criação de uma fatura
func ContinuaCriacaoFatura(bot *tgbotapi.BotAPI, message tgbotapi.Update, userState *UserState) {
	switch userState.CurrentStep {
	case "cadastro_fatura":
		edit := tgbotapi.NewEditMessageText(message.CallbackQuery.Message.Chat.ID, message.CallbackQuery.Message.MessageID, "Cartão selecionado 😀")
		edit.ReplyMarkup = nil

		_, err := bot.Send(edit)
		if err != nil {
			log.Panic(err)
		}

		userState.NewInvoiceData.CartaoID = message.CallbackQuery.Data
		userState.CurrentStep = "cadastro_fatura_titulo"
		msg := tgbotapi.NewMessage(message.CallbackQuery.Message.Chat.ID, "Por favor, insira a data de vencimento da nova fatura, ex: 20-02-2024")
		_, err = bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	default:
		enviaErroMensagemInformadaEstadoDesconhecido(bot, message.CallbackQuery.Message.Chat.ID)
	}
}

// finalizarCadastroFatura é responsável por finalizar o processo de criação de uma fatura
func finalizarCadastroFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userTokens map[int64]string, userState *UserState) {
	if userState.CurrentStep == "cadastro_fatura_titulo" {
		data, _ := time.Parse("02-01-2006", message.Text)
		dataFormatada := data.Format(time.DateOnly)
		userState.NewInvoiceData.DataVencimento = &dataFormatada

		err := CadastrarFatura(userState, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		userState.CurrentStep = ""
		userState.CurrentStepBool = false
		msg := tgbotapi.NewMessage(message.Chat.ID, "Fatura cadastrada com sucesso")
		utils.EnviaMensagem(bot, msg)
	} else {
		enviaErroMensagemInformadaEstadoDesconhecido(bot, message.Chat.ID)
	}
}

// enviaErroMensagemInformadaEstadoDesconhecido é uma Função para enviar uma mensagem informando um erro de estado desconhecido
func enviaErroMensagemInformadaEstadoDesconhecido(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Erro: Estado desconhecido.")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func gerarOpcoesFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Criando um teclado de resposta
	buttonOpcao1 := tgbotapi.NewKeyboardButton("Compras Fatura")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("Cadastrar Fatura")
	buttonOpcao3 := tgbotapi.NewKeyboardButton("Pagar Fatura")
	buttonOpcao4 := tgbotapi.NewKeyboardButton("Opção 4")

	keyboard := tgbotapi.NewReplyKeyboard(
		[]tgbotapi.KeyboardButton{buttonOpcao1, buttonOpcao2},
		[]tgbotapi.KeyboardButton{buttonOpcao3, buttonOpcao4},
	)

	// Configurando a mensagem de boas-vindas com o teclado de resposta
	msg := tgbotapi.NewMessage(message.Chat.ID, "Selecione uma opção:")
	msg.ReplyMarkup = keyboard

	// Enviando a mensagem
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// EnviarOpcoesCartoesFatura Função para enviar botões inline de seleção de cartões
func EnviarOpcoesCartoesFatura(bot *tgbotapi.BotAPI, chatID int64, cartao *cartao.ResPag, userStates *UserStepComprasFatura) {
	// Criar slice para armazenar botões
	var buttons [][]tgbotapi.InlineKeyboardButton

	for i, card := range cartao.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*card.Nome, card.ID.String())

		// Adicionar botão à linha atual
		row := i / 3
		if len(buttons) <= row {
			// Adicionar uma nova linha se necessário
			buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})
		}
		buttons[row] = append(buttons[row], button)
	}

	// Criar teclado inline com os botões
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Configurar a mensagem com o teclado
	msg := tgbotapi.NewMessage(chatID, "Selecione um cartão: ")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "fatura_selecionada"

	userStates.Opcao = &step
}

// EnviarOpcoesFaturas Função para enviar botões inline de seleção de faturas
func EnviarOpcoesFaturas(bot *tgbotapi.BotAPI, chatID int64, faturas *ResPagFaturas, userStates *UserStepComprasFatura, userCompras *compras.UserStateCompras, callbackQuery *tgbotapi.CallbackQuery, userTokens map[int64]string) {
	res, err := cartao.BuscarCartao(cartao.BaseURLCartao, callbackQuery.Data, userTokens, callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Cartão Selecionado: %s", *res.Nome))
	edit.ReplyMarkup = nil

	_, err = bot.Send(edit)
	if err != nil {
		log.Panic(err)
	}

	var buttons [][]tgbotapi.InlineKeyboardButton

	// Adicionar botões para cada fatura
	for i, invoice := range faturas.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*invoice.Nome, invoice.ID.String())

		// Adicionar botão à linha atual
		row := i / 3
		if len(buttons) <= row {
			// Adicionar uma nova linha se necessário
			buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})
		}
		buttons[row] = append(buttons[row], button)
	}

	// Criar teclado inline com os botões
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	// Configurar a mensagem com o teclado
	msg := tgbotapi.NewMessage(chatID, "Selecione uma fatura:")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	var step string

	if userStates.Opcao != nil {
		step = "cartao_fatura_selecionado"

		userStates.Opcao = &step
	} else {
		step = "fatura_selecionada"
		*userCompras.CurrentStep = step
	}
}

// ProcessarCasosStepComprasFatura é responsável por controlar o fluxo que obtém as compras de uma fatura
func ProcessarCasosStepComprasFatura(userCompraFaturas *UserStepComprasFatura, reqStatus *ReqAtualizarStatus, bot *tgbotapi.BotAPI, update tgbotapi.Update, userTokens map[int64]string) {
	switch *userCompraFaturas.Opcao {
	case "fatura_selecionada":
		faturasCartao, err := ListarFaturas(fmt.Sprintf(BaseURLFaturas+"%s/faturas?status=Em Aberto", update.CallbackQuery.Data), userTokens, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, nil, update.CallbackQuery, userTokens)
	case "cartao_fatura_selecionado":
		ProcessCallbackQuery(bot, update.CallbackQuery, userTokens)
	case "Pagar Fatura":
		faturasCartao, err := ListarFaturas(fmt.Sprintf(BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data), userTokens, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, nil, update.CallbackQuery, userTokens)

		*userCompraFaturas.Opcao = "pagar_fatura_selecionado"
	case "pagar_fatura_selecionado":
		faturaID, err := uuid.Parse(update.CallbackQuery.Data)
		if err != nil {
			log.Panic(err)
		}

		userCompraFaturas.Fatura.ID = &faturaID

		EnviaOpcoesNovoStatusFatura(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userCompraFaturas, userTokens)
	case "status_selecionado":
		reqStatus.Status = &update.CallbackQuery.Data

		err := AtualizarStatusPagamentoFatura(userCompraFaturas.Fatura.ID, reqStatus, userTokens, update.CallbackQuery.Message.Chat.ID)
		if err != nil {
			edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, err.Error())
			utils.EnviaMensagem(bot, edit)
			return
		}

		edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Novo Status da Fatura: %s", update.CallbackQuery.Data))
		edit.ReplyMarkup = nil
		utils.EnviaMensagem(bot, edit)
	}
}

// ProcessCallbackQuery Função respnsável para processar a escolha do usuário, mostrando as compras realizadas e gerando um pdf com elas
func ProcessCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, userTokens map[int64]string) {
	res, err := BuscarFatura(&callbackQuery.Data, userTokens, callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	pdfContent, err := compras.ObterComprasPdf(res.ID, res.FaturaCartaoID, userTokens, callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	msgPdfCompras := tgbotapi.NewDocumentUpload(callbackQuery.Message.Chat.ID, tgbotapi.FileReader{
		Name:   "compras_" + strings.ToLower(*res.Nome) + "_" + strings.ToLower(*res.NomeCartao) + ".pdf",
		Reader: bytes.NewBuffer(pdfContent),
		Size:   int64(len(pdfContent)),
	})

	// Realizar ações com base no ID da fatura selecionada
	log.Printf("Usuário selecionou a fatura com ID: %s", *res.Nome)

	// Responder para indicar que a callback query foi processada
	answer := tgbotapi.NewCallback(callbackQuery.ID, "")
	_, err = bot.AnswerCallbackQuery(answer)
	if err != nil {
		log.Println("Erro ao responder à callback query:", err)
		return
	}

	dataVencimentoFormat := *res.DataVencimento

	// Opcional: editar a mensagem original para remover os botões
	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID,
		fmt.Sprintf("Dados básicos da Fatura selecionada: \n\n"+
			"Nome Cartão: %s \n"+
			"Status Pagamento: %s \n"+
			"Data Vencimento: %s \n", *res.NomeCartao, *res.Status, dataVencimentoFormat[:10]))
	edit.ReplyMarkup = nil // Remove o teclado inline

	comprasFatura, err := compras.ListarComprasFatura(&callbackQuery.Data, userTokens, callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	var (
		msgRetornoCompras string
	)

	for _, compra := range comprasFatura.Dados {
		dataCompraFormat := *compra.DataCompra
		msgRetornoCompra := fmt.Sprintf("Local da Compra: %s \n"+
			"Descrição: %s \n"+
			"Categoria: %s \n"+
			"Valor da Parcela: %.2f \n"+
			"Parcela Atual: %d \n"+
			"Quantidade de Parcelas: %d \n"+
			"Data da Compra: %s \n\n", *compra.LocalCompra, *compra.Descricao, *compra.CategoriaNome, *compra.ValorParcela, *compra.ParcelaAtual, *compra.QuantidadeParcelas, dataCompraFormat[:10])

		msgRetornoCompras += msgRetornoCompra
	}

	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, msgRetornoCompras)

	utils.EnviaMensagem(bot, msg)
	utils.EnviaMensagem(bot, edit)
	utils.EnviaMensagem(bot, msgPdfCompras)
}

// EnviaOpcoesNovoStatusFatura realiza o envio das opções disponíveis para mudar o status de uma fatura
func EnviaOpcoesNovoStatusFatura(bot *tgbotapi.BotAPI, chatID int64, callbackQuery *tgbotapi.CallbackQuery, userCompraFaturas *UserStepComprasFatura, userTokens map[int64]string) {
	var (
		buttons []tgbotapi.InlineKeyboardButton
		options = []string{"Em Aberto", "Pago", "Atrasada"}
	)

	res, err := BuscarFatura(&callbackQuery.Data, userTokens, callbackQuery.Message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, err.Error())
		utils.EnviaMensagem(bot, msg)
		return
	}

	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Status atual da fatura: %s", *res.Status))
	edit.ReplyMarkup = nil
	utils.EnviaMensagem(bot, edit)

	for _, status := range options {
		button := tgbotapi.NewInlineKeyboardButtonData(status, status)

		buttons = append(buttons, button)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)

	msg := tgbotapi.NewMessage(chatID, "Selecione o novo status da fatura!")
	msg.ReplyMarkup = keyboard

	// Enviar a mensagem
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "status_selecionado"

	userCompraFaturas.Opcao = &step
}

func ListarFaturas(url string, userTokens map[int64]string, chatID int64) (res ResPagFaturas, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return res, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+url, nil)
	if err != nil {
		return res, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return res, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

// BuscarFatura é responsável por realizar uma requisição para obter os dados de uma fatura
func BuscarFatura(id *string, userTokens map[int64]string, chatID int64) (res Res, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return res, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+BaseURLFatura+*id, nil)
	if err != nil {
		return res, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return res, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("%s", resp.Status)
	}

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err = json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

// AtualizarStatusPagamentoFatura é responsável por realizar a requisição para atualização do status de uma fatura
func AtualizarStatusPagamentoFatura(faturaID *uuid.UUID, dadosStatus *ReqAtualizarStatus, userTokens map[int64]string, chatID int64) (err error) {
	dados := ReqAtualizarStatus{
		Status: dadosStatus.Status,
	}

	token, ok := userTokens[chatID]
	if !ok {
		return fmt.Errorf("usuário não está autenticado")
	}

	dadosJSON, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao codificar os dados JSON:", err)
		return
	}

	var ambiente = utils.ValidarAmbiente()

	req, err := http.NewRequest(http.MethodPut, ambiente+BaseURLFatura+fmt.Sprintf("%s/status", faturaID.String()), bytes.NewBuffer(dadosJSON))
	if err != nil {
		fmt.Println("Erro ao criar a requisição PUT:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição PUT:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s", resp.Status)
	}

	return
}

// CadastrarFatura é responsável por realizar a requisição para cadastrar uma fatura
func CadastrarFatura(fatura *UserState, userTokens map[int64]string, chatID int64) (err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return fmt.Errorf("usuário não está autenticado")
	}

	var baseURLCadastroFaturas = fmt.Sprintf("%s/cadastros/cartao", ambiente)
	// Montar os dados a serem enviados no corpo do POST
	dados := NewInvoice{
		CartaoID:       fatura.NewInvoiceData.CartaoID,
		DataVencimento: fatura.NewInvoiceData.DataVencimento,
	}

	// Codificar os dados em formato JSON
	dadosJSON, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao codificar os dados JSON:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, baseURLCadastroFaturas+fmt.Sprintf("/%s/faturas", fatura.NewInvoiceData.CartaoID), bytes.NewBuffer(dadosJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Realize login!!")
	} else if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro ao realizar cadastro da fatura")
	}

	return
}
