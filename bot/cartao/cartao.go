package cartao

import (
	"bot_controle_cartao/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ProcessoAcoesCartoes é responsável por coordenar as ações relacionadas a cartões
func ProcessoAcoesCartoes(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCartaoState map[int64]*UserStateCartao, userTokens map[int64]string) {
	if userState, ok := userCartaoState[message.Chat.ID]; ok {
		if userState.CurrentStepBool {
			continuaCriacaoCartao(bot, message, userState, userTokens)
		}
	}

	switch message.Text {
	case "Cartões":
		gerarOpcoesAcoesCartao(bot, message)
	case "Extrato":
		cartoes, err := ListarCartoes(BaseURLCartoes, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		gerarOpcoesCartoesDisponiveis(bot, message.Chat.ID, &cartoes, userCartaoState)
	case "Cadastrar Cartão":
		inicioCriacaoCartao(bot, message.Chat.ID, userCartaoState)
		userCartaoState[message.Chat.ID].CurrentStepBool = true
	}
}

// inicioCriacaoCartao é responsável por iniciar o processo de criação de um cartão
func inicioCriacaoCartao(bot *tgbotapi.BotAPI, chatID int64, userCartaoState map[int64]*UserStateCartao) {
	userCartaoState[chatID] = &UserStateCartao{
		ChatID:      chatID,
		CurrentStep: "cadastro_cartao",
	}

	msg := tgbotapi.NewMessage(chatID, "Por favor, insira o nome do cartão")
	utils.EnviaMensagem(bot, msg)
}

// continuaCriacaoCartao é responsável por continuar a criação do cartão
func continuaCriacaoCartao(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCartaoState *UserStateCartao, userTokens map[int64]string) {
	switch userCartaoState.CurrentStep {
	case "cadastro_cartao":
		userCartaoState.NovoCartaoData.Nome = message.Text
		userCartaoState.CurrentStep = ""
		userCartaoState.CurrentStepBool = false

		err := CadastrarCartao(userCartaoState, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Cartão cadastrado com sucesso!")
		utils.EnviaMensagem(bot, msg)
	}
}

// CadastrarCartao é responsável por cadastrar um novo cartão
func CadastrarCartao(cartao *UserStateCartao, userTokens map[int64]string, chatID int64) (err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return fmt.Errorf("usuário não está autenticado")
	}

	// Montar os dados a serem enviados no corpo do POST
	dados := NovoCartao{
		Nome: cartao.NovoCartaoData.Nome,
	}

	// Codificar os dados em formato JSON
	dadosJSON, err := json.Marshal(dados)
	if err != nil {
		fmt.Println("Erro ao codificar os dados JSON:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s"+BaseURLCartoes, ambiente), bytes.NewBuffer(dadosJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Erro ao realizar cadastro do cartão!")
	}

	return
}

// gerarOpcoesAcoesCartao é responsável por gerar os botões para seleção das ações de cartões para o usuário
func gerarOpcoesAcoesCartao(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	buttonOpcao1 := tgbotapi.NewKeyboardButton("Extrato")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("Cadastrar Cartão")

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

// gerarOpcoesCartoesDisponiveis Função para enviar botões inline de seleção de cartões
func gerarOpcoesCartoesDisponiveis(bot *tgbotapi.BotAPI, chatID int64, cartao *ResPag, userCartaoState map[int64]*UserStateCartao) {

	var buttons [][]tgbotapi.InlineKeyboardButton

	// Adicionar botões para cada fatura
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

	step := "selecionar_ano"

	userCartaoState[chatID] = &UserStateCartao{
		ChatID:      chatID,
		CurrentStep: step,
	}
}

// EnviarOpcoesAno envia as opções para que o usuário selecione o ano que será usado quando for gerado o extrato
func EnviarOpcoesAno(bot *tgbotapi.BotAPI, chatID int64, callbackQuery *tgbotapi.CallbackQuery, userCartaoState map[int64]*UserStateCartao) {
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

	userCartaoState[chatID].CurrentStep = step
}

// ListarCartoes é responsável por listar os cartões cadastrados
func ListarCartoes(url string, userTokens map[int64]string, chatID int64) (cartoes ResPag, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return cartoes, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+url, nil)
	if err != nil {
		return cartoes, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return cartoes, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return cartoes, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err = json.Unmarshal(body, &cartoes); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

// BuscarCartao é responsável por buscar o cartão de acordo com o id
func BuscarCartao(url string, id string, userTokens map[int64]string, chatID int64) (cartao Res, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return cartao, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+url+fmt.Sprintf("/%s", id), nil)
	if err != nil {
		return cartao, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return cartao, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return cartao, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err = json.Unmarshal(body, &cartao); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}
