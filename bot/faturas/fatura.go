package faturas

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/compras"
	"bot_controle_cartao/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ProcessoAcoesFaturas é responsável por ter todos os processos das faturas
func ProcessoAcoesFaturas(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userStates map[int64]*UserState, userCompraFatura *UserStepComprasFatura) {
	// Verificar o texto da mensagem para determinar a ação a ser tomada

	if userState, ok := userStates[message.Chat.ID]; ok {
		if userState.CurrentStepBool {
			continuaCriacaoFatura(bot, message, userState)
		}
	}

	switch message.Text {
	case "Faturas":
		gerarOpcoesFatura(bot, message)
	case "Compras Fatura":
		cartoes := cartao.ListarCartoes(cartao.BaseURLCartoes)

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCompraFatura)
	case "cadastrar_fatura":
		inicioCriacaoFatura(bot, message.Chat.ID, userStates)
		userStates[message.Chat.ID].CurrentStepBool = true
	case "Pagar Fatura":
		cartoes := cartao.ListarCartoes(cartao.BaseURLCartoes)

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCompraFatura)

		*userCompraFatura.Opcao = "Pagar Fatura"
	}
}

// EnviaMensagemBoasVindas é responsável por enviar uma mensagem de boas vindas
func EnviaMensagemBoasVindas(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Bem-vindo ao Bot de Faturas!")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// inicioCriacaoFatura é responsável por iniciar o processo de criação de uma fatura
func inicioCriacaoFatura(bot *tgbotapi.BotAPI, chatID int64, userStates map[int64]*UserState) {
	// Definir o estado da conversa do usuário como "cadastro_fatura"
	userStates[chatID] = &UserState{
		ChatID:      chatID,
		CurrentStep: "cadastro_fatura",
	}

	msg := tgbotapi.NewMessage(chatID, "Por favor, insira o título da fatura:")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// continuaCriacaoFatura é responsável por continuar o processo de criação de uma fatura
func continuaCriacaoFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *UserState) {
	switch userState.CurrentStep {
	case "cadastro_fatura":
		// Armazenar o título da fatura
		userState.NewInvoiceData.Title = message.Text
		userState.CurrentStep = "cadastro_fatura_valor"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, insira o valor da fatura:")
		_, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	case "cadastro_fatura_valor":
		// Armazenar o valor da fatura, deve fazer a validação do valor
		amountValue, erro := strconv.ParseFloat(message.Text, 64)
		if erro != nil {
			log.Panic(erro)
		}
		userState.NewInvoiceData.Amount = amountValue
		userState.CurrentStep = "cadastro_fatura_data_vencimento"
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, insira a data de vencimento da fatura:")
		_, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	case "cadastro_fatura_data_vencimento":
		userState.NewInvoiceData.DueDate = message.Text
		userState.CurrentStep = "" // Resetar o estado da conversa do usuário
		userState.CurrentStepBool = false
		// Aqui você pode processar os dados da nova fatura, e enviar uma mensagem de confirmação ao usuário
		enviaMensagemCadastroFaturaSucesso(bot, message.Chat.ID, userState.NewInvoiceData)
	default:
		// Se o estado da conversa do usuário não for reconhecido, enviar uma mensagem informando o erro
		enviaErroMensagemInformadaEstadoDesconhecido(bot, message.Chat.ID)
	}
}

// enviaMensagemCadastroFaturaSucesso é uma Função para enviar uma mensagem de confirmação com os detalhes da fatura cadastrada
func enviaMensagemCadastroFaturaSucesso(bot *tgbotapi.BotAPI, chatID int64, newInvoice NewInvoice) {
	msgText := "Nova fatura cadastrada com sucesso!\n\n" +
		"Título: " + newInvoice.Title + "\n" +
		"Valor: " + "R$ " + FormataValorFatura(newInvoice.Amount) + "\n" +
		"Data de Vencimento: " + newInvoice.DueDate

	msg := tgbotapi.NewMessage(chatID, msgText)
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

// FormataValorFatura Função para formatar o valor da fatura
func FormataValorFatura(amount float64) string {
	return strconv.FormatFloat(amount, 'f', -1, 64)
}

// enviaErroMensagemComandoNaoReconhecido é uma Função para enviar uma mensagem informando que o comando não foi reconhecido
func enviaErroMensagemComandoNaoReconhecido(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Comando não reconhecido.")
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
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
	buttonOpcao2 := tgbotapi.NewKeyboardButton("cadastrar_fatura")
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
func EnviarOpcoesFaturas(bot *tgbotapi.BotAPI, chatID int64, faturas *ResPagFaturas, userStates *UserStepComprasFatura, userCompras *compras.UserStateCompras, callbackQuery *tgbotapi.CallbackQuery) {
	res := cartao.BuscarCartao(cartao.BaseURLCartao, callbackQuery.Data)

	edit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("Cartão Selecionado: %s", *res.Nome))
	edit.ReplyMarkup = nil

	_, err := bot.Send(edit)
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
func ProcessarCasosStepComprasFatura(userCompraFaturas *UserStepComprasFatura, reqStatus *ReqAtualizarStatus, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch *userCompraFaturas.Opcao {
	case "fatura_selecionada":
		faturasCartao := ListarFaturas(fmt.Sprintf(BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data))

		EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, nil, update.CallbackQuery)
	case "cartao_fatura_selecionado":
		ProcessCallbackQuery(bot, update.CallbackQuery)
	case "Pagar Fatura":
		faturasCartao := ListarFaturas(fmt.Sprintf(BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data))

		EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, nil, update.CallbackQuery)

		*userCompraFaturas.Opcao = "pagar_fatura_selecionado"
	case "pagar_fatura_selecionado":
		faturaID, err := uuid.Parse(update.CallbackQuery.Data)
		if err != nil {
			log.Panic(err)
		}

		userCompraFaturas.Fatura.ID = &faturaID

		EnviaOpcoesNovoStatusFatura(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userCompraFaturas)
	case "status_selecionado":
		reqStatus.Status = &update.CallbackQuery.Data

		err := AtualizarStatusPagamentoFatura(userCompraFaturas.Fatura.ID, reqStatus)
		if err.Err != nil {
			edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, *err.Err)
			utils.EnviaMensagem(bot, edit)
			return
		}

		edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Novo Status da Fatura: %s", update.CallbackQuery.Data))
		edit.ReplyMarkup = nil
		utils.EnviaMensagem(bot, edit)
	}
}

// ProcessCallbackQuery Função respnsável para processar a escolha do usuário, mostrando as compras realizadas e gerando um pdf com elas
func ProcessCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Obter o ID da fatura a partir do CallbackData
	res := BuscarFatura(&callbackQuery.Data)

	pdfContent := compras.ObterComprasPdf(res.ID, res.FaturaCartaoID)

	msgPdfCompras := tgbotapi.NewDocumentUpload(callbackQuery.Message.Chat.ID, tgbotapi.FileReader{
		Name:   "compras_" + strings.ToLower(*res.Nome) + "_" + strings.ToLower(*res.NomeCartao) + ".pdf",
		Reader: bytes.NewBuffer(pdfContent),
		Size:   int64(len(pdfContent)),
	})

	// Realizar ações com base no ID da fatura selecionada
	log.Printf("Usuário selecionou a fatura com ID: %s", *res.Nome)

	// Responder para indicar que a callback query foi processada
	answer := tgbotapi.NewCallback(callbackQuery.ID, "")
	_, err := bot.AnswerCallbackQuery(answer)
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

	comprasFatura := compras.ListarComprasFatura(&callbackQuery.Data)

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
func EnviaOpcoesNovoStatusFatura(bot *tgbotapi.BotAPI, chatID int64, callbackQuery *tgbotapi.CallbackQuery, userCompraFaturas *UserStepComprasFatura) {
	res := BuscarFatura(&callbackQuery.Data)
	var (
		buttons []tgbotapi.InlineKeyboardButton
		options = []string{"Em Aberto", "Pago", "Atrasada"}
	)

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
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	step := "status_selecionado"

	userCompraFaturas.Opcao = &step
}

func ListarFaturas(url string) (res ResPagFaturas) {
	var ambiente = utils.ValidarAmbiente()

	resp, err := http.Get(ambiente + url)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
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

// BuscarFatura é responsável por realizar uma requisição para obter os dados de uma fatura
func BuscarFatura(id *string) (res Res) {
	var ambiente = utils.ValidarAmbiente()

	resp, err := http.Get(ambiente + BaseURLFatura + *id)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

type Error struct {
	Err *string `json:"error"`
}

// AtualizarStatusPagamentoFatura é responsável por realizar a requisição para atualização do status de uma fatura
func AtualizarStatusPagamentoFatura(faturaID *uuid.UUID, dadosStatus *ReqAtualizarStatus) (estruturaErro Error) {
	dados := ReqAtualizarStatus{
		Status: dadosStatus.Status,
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição PUT:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	//if resp.StatusCode != http.StatusOK {
	//	err = utils.NewErr(string(body))
	//	return err
	//}

	// Imprime a resposta da API
	fmt.Println("Resposta da API:", string(body))

	if err := json.Unmarshal(body, &estruturaErro); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}
