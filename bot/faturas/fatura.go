package faturas

import (
	"bot_controle_cartao/cartao"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	case "faturas":
		gerarOpcoesFatura(bot, message)
	case "/compras":
		inicioComprasFatura(message.Chat.ID, userStates)
		userCompraFatura.ComprasFatura = true
	case "/cadastrar_fatura":
		inicioCriacaoFatura(bot, message.Chat.ID, userStates)
		userStates[message.Chat.ID].CurrentStepBool = true
	}

	if userCompraFatura.ComprasFatura {
		if userState, ok := userStates[message.Chat.ID]; ok {
			continuaComprasFatura(bot, message, userState, userCompraFatura)
		}
	}

}

// inicioComprasFatura é responsável por iniciar o processo de obter as compras de uma fatura
func inicioComprasFatura(chatID int64, userStates map[int64]*UserState) {
	// Definir o estado da conversa do usuário como "cadastro_fatura"
	userStates[chatID] = &UserState{
		ChatID:      chatID,
		CurrentStep: "start_compras_fatura",
	}
}

// enviaMensagemInicio é responsável por enviar uma mensagem de boas vindas
func enviaMensagemInicio(bot *tgbotapi.BotAPI, chatID int64) {
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

// continuaComprasFatura é responsável por continuar o processo de solicitação das compras de uma fatura
func continuaComprasFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *UserState, userCompraFatura *UserStepComprasFatura) {
	switch userState.CurrentStep {
	case "start_compras_fatura":
		obterCartoesFatura(bot, message, userState, userCompraFatura)
	case "compras_fatura":
		obterCartaoSelecionado(bot, message, userState, userCompraFatura)
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
	buttonOpcao1 := tgbotapi.NewKeyboardButton("/compras")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("/cadastrar_fatura")
	buttonOpcao3 := tgbotapi.NewKeyboardButton("Opção 3")
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

func obterCartoesFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *UserState, userCompraFatura *UserStepComprasFatura) {
	cartoes := cartao.ListarCartoes(BaseURLCartoes)

	var (
		options []string
	)
	for _, cartaoValue := range cartoes.Dados {
		options = append(options, *cartaoValue.Nome)
	}

	userCompraFatura.Cartoes = options

	var buttons []tgbotapi.KeyboardButton

	for i := range options {
		button := tgbotapi.NewKeyboardButton(options[i])
		buttons = append(buttons, button)
	}

	lenOptions := (len(buttons) + 1) / 2

	buttonsOne := make([]tgbotapi.KeyboardButton, lenOptions)
	buttonsTwo := make([]tgbotapi.KeyboardButton, lenOptions)

	copy(buttonsOne, buttons[:lenOptions])
	copy(buttonsTwo, buttons[lenOptions:])

	keyboard := tgbotapi.NewReplyKeyboard(
		buttonsOne,
		buttonsTwo,
	)

	// Configurando a mensagem de boas-vindas com o teclado de resposta
	msg := tgbotapi.NewMessage(message.Chat.ID, "Selecione o cartão:")
	msg.ReplyMarkup = keyboard

	// Enviando a mensagem
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	userState.CurrentStep = "compras_fatura"

	return
}

func obterCartaoSelecionado(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userState *UserState, userCompraFatura *UserStepComprasFatura) {
	if len(userCompraFatura.Cartoes) != 0 {

		for _, option := range userCompraFatura.Cartoes {
			if message.Text == option {
				userCompraFatura.Opcao = &option
				break
			}
		}
	}

	obterMesesFatura(bot, message, userCompraFatura)

	return
}

func obterMesesFatura(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCompraFatura *UserStepComprasFatura) {
	var options []string

	if userCompraFatura.Opcao != nil {
		cartaoEspecifico := cartao.ListarCartoes(fmt.Sprintf(BaseURLCartoes+"?nome_exato=%v", *userCompraFatura.Opcao))

		faturas := ListarFaturas(fmt.Sprintf(BaseURLFaturas+"%v/faturas", cartaoEspecifico.Dados[0].ID))

		for _, fatura := range faturas.Dados {
			options = append(options, *fatura.Nome)
		}

		var buttonsMesesFatura []tgbotapi.KeyboardButton

		for i := range options {
			button := tgbotapi.NewKeyboardButton(options[i])
			buttonsMesesFatura = append(buttonsMesesFatura, button)
		}

		lenFaturas := (len(buttonsMesesFatura) + 1) / 2

		buttonsOneFaturas := make([]tgbotapi.KeyboardButton, lenFaturas)
		buttonsTwoFaturas := make([]tgbotapi.KeyboardButton, lenFaturas)

		copy(buttonsOneFaturas, buttonsMesesFatura[:lenFaturas])
		copy(buttonsTwoFaturas, buttonsMesesFatura[lenFaturas:])

		keyboardFaturas := tgbotapi.NewReplyKeyboard(
			buttonsOneFaturas,
			buttonsTwoFaturas,
		)

		// Configurando a mensagem de boas-vindas com o teclado de resposta
		msgFaturas := tgbotapi.NewMessage(message.Chat.ID, "Selecione a fatura:")
		msgFaturas.ReplyMarkup = keyboardFaturas

		// Enviando a mensagem
		_, err := bot.Send(msgFaturas)
		if err != nil {
			log.Panic(err)
		}
	}

	userCompraFatura.ComprasFatura = false
}

func ListarFaturas(url string) (res ResPagFaturas) {
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

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}
