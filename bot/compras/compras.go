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
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ProcessoAcoesCompras é responsável por coordenar as ações relacionadas a compras
func ProcessoAcoesCompras(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCartaoState *UserStateCompras, acaoAnterior string, userTokens map[int64]string) {

	switch message.Text {
	case "Compras":
		gerarOpcoesAcoesCompras(bot, message)
	case "Cadastrar Compra":
		cartoes, err := cartao.ListarCartoes(cartao.BaseURLCartoes, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		EnviarOpcoesCartoesFatura(bot, message.Chat.ID, &cartoes, userCartaoState)
	case "Obter total compras":
		InicioObterTotalCompras(userCartaoState)
	}

	if acaoAnterior == "compras" {
		ProcessamentoObterTotalCompras(bot, message, userCartaoState, userTokens)
	}
}

// ProcessamentoObterTotalCompras é responsável por processar o fluxo para obter o valor total compras
func ProcessamentoObterTotalCompras(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCompras *UserStateCompras, userTokens map[int64]string) {
	var cancelado bool

	switch *userCompras.CurrentStep {
	case "selecionar_data":
		*userCompras.CurrentStep = "selecionar_pago"

		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Digite a data esperada para filtragem. Ex: 08/2024 \n")
		utils.EnviaMensagem(bot, msg)

		msg = tgbotapi.NewMessage(message.Chat.ID, "A qualquer momento digite 'cancelar' para cancelar essa operação \n")
		utils.EnviaMensagem(bot, msg)
	case "selecionar_pago":
		cancelado = utils.CancelarOperacao(bot, &message.Text, userCompras.CurrentStep, message.Chat.ID)

		if !cancelado {
			*userCompras.CurrentStep = "selecionar_ultima_parcela"

			userCompras.ObterTotalCompras.DataEspecifica = &message.Text

			msg := tgbotapi.NewMessage(message.Chat.ID, "Você gostaria que fosse excluído da soma, os valores que já foram pagos? Digite: SIM ou NÃO")
			utils.EnviaMensagem(bot, msg)
		}
	case "selecionar_ultima_parcela":
		cancelado = utils.CancelarOperacao(bot, &message.Text, userCompras.CurrentStep, message.Chat.ID)

		if !cancelado {
			switch strings.ToLower(message.Text) {
			case "sim":
				userCompras.ObterTotalCompras.Pago = utils.ObterPonteiro[bool](true)
				*userCompras.CurrentStep = "requisicao_obter_valor"

				msg := tgbotapi.NewMessage(message.Chat.ID, "Você gostaria que retornasse a soma apenas dos valores que estão na última parcela? Digite: SIM ou NÃO")
				utils.EnviaMensagem(bot, msg)
			case "não":
				userCompras.ObterTotalCompras.Pago = utils.ObterPonteiro[bool](false)
				*userCompras.CurrentStep = "requisicao_obter_valor"

				msg := tgbotapi.NewMessage(message.Chat.ID, "Você gostaria que retornasse a soma apenas dos valores que estão na última parcela? Digite: SIM ou NÃO")
				utils.EnviaMensagem(bot, msg)
			default:
				msg := tgbotapi.NewMessage(message.Chat.ID, "Valor incorreto, digite novamente")
				utils.EnviaMensagem(bot, msg)
			}
		}
	case "requisicao_obter_valor":
		cancelado = utils.CancelarOperacao(bot, &message.Text, userCompras.CurrentStep, message.Chat.ID)

		if !cancelado {
			switch strings.ToLower(message.Text) {
			case "sim":
				userCompras.ObterTotalCompras.UltimaParcela = utils.ObterPonteiro[bool](true)
				userCompras.CurrentStep = nil

				res, err := ObterComprasTotal(userCompras, userTokens, message.Chat.ID)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
					utils.EnviaMensagem(bot, msg)
					return
				}

				msg := tgbotapi.NewMessage(message.Chat.ID, "Valor total: "+*res.Total)
				utils.EnviaMensagem(bot, msg)
			case "não":
				userCompras.ObterTotalCompras.UltimaParcela = utils.ObterPonteiro[bool](false)
				userCompras.CurrentStep = nil

				res, err := ObterComprasTotal(userCompras, userTokens, message.Chat.ID)
				if err != nil {
					msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
					utils.EnviaMensagem(bot, msg)
					return
				}

				msg := tgbotapi.NewMessage(message.Chat.ID, "Valor total: "+*res.Total)
				utils.EnviaMensagem(bot, msg)
			default:
				msg := tgbotapi.NewMessage(message.Chat.ID, "Valor incorreto, digite novamente")
				utils.EnviaMensagem(bot, msg)
			}
		}
	}
}

// ProcessoAcoesCadastroCompra é responsável pelo fluxo de cadastro de uma compra
func ProcessoAcoesCadastroCompra(bot *tgbotapi.BotAPI, message *tgbotapi.Message, userCompras *UserStateCompras, userTokens map[int64]string) {
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
		msg := tgbotapi.NewMessage(message.Chat.ID, "Por favor, Diga qual foi a data da compra, ex: 20-02-2024")

		utils.EnviaMensagem(bot, msg)
	case "quantidade_parcelas_selecionada":
		data, _ := time.Parse("02-01-2006", message.Text)
		dataFormatada := data.Format(time.DateOnly)

		userCompras.NovaCompraData.DataCompra = &dataFormatada
		*userCompras.CurrentStep = ""

		err := CadastrarCompra(userCompras, userTokens, message.Chat.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			utils.EnviaMensagem(bot, msg)
			return
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Compra cadastrada com sucesso!")
		utils.EnviaMensagem(bot, msg)
	}
}

// gerarOpcoesAcoesCompras é responsável por gerar os botões para seleção das ações de compras para o usuário
func gerarOpcoesAcoesCompras(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	buttonOpcao1 := tgbotapi.NewKeyboardButton("Cadastrar Compra")
	buttonOpcao2 := tgbotapi.NewKeyboardButton("Obter total compras")
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

	step := "selecionar_fatura"

	userStatesCompras.CurrentStep = &step
}

// InicioObterTotalCompras é responsável por enviar as opções de filtragem
func InicioObterTotalCompras(stateCompras *UserStateCompras) {
	step := "selecionar_data"
	stateCompras.CurrentStep = &step
	stateCompras.ObterTotalCompras.StepComprasTotal = utils.ObterPonteiro(true)
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
	var buttons [][]tgbotapi.InlineKeyboardButton

	for i, categoria := range categorias.Dados {
		button := tgbotapi.NewInlineKeyboardButtonData(*categoria.Nome, categoria.ID.String())

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
func ListarComprasFatura(idFatura *string, userTokens map[int64]string, chatID int64) (res ResComprasPag, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return res, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+BaseURLCompras+"?fatura_id="+*idFatura, nil)
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

// CadastrarCompra é responsável por realizar a requisição para cadastrar uma compra
func CadastrarCompra(compras *UserStateCompras, userTokens map[int64]string, chatID int64) (err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return fmt.Errorf("usuário não está autenticado")
	}

	var baseURLCadastroCompras = fmt.Sprintf("%s/cadastros/fatura", ambiente)
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

	req, err := http.NewRequest(http.MethodPost, baseURLCadastroCompras+fmt.Sprintf("/%s/compras", *compras.FaturaID), bytes.NewBuffer(dadosJSON))
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
		return fmt.Errorf("Erro ao realizar cadastro da compra!")
	}

	return
}

// ObterComprasPdf é responsável por realizar a requisição que obtém o pdf com as compras
func ObterComprasPdf(idFatura *uuid.UUID, idCartao *uuid.UUID, userTokens map[int64]string, chatID int64) (body []byte, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)

	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return nil, fmt.Errorf("usuário não está autenticado")
	}

	if idFatura != nil && idCartao != nil {
		req, err = http.NewRequest(http.MethodGet, ambiente+BaseURLComprasPdf+"?fatura_id="+idFatura.String()+"&cartao_id="+idCartao.String(), nil)
		if err != nil {
			return
		}
	} else {
		req, err = http.NewRequest(http.MethodGet, ambiente+BaseURLComprasPdf+"?cartao_id="+idCartao.String(), nil)
		if err != nil {
			return
		}
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return body, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("%s", resp.Status)
	}

	// Lê o corpo da resposta
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	return
}

// ObterComprasTotal é responsável por realizar a requisição que obtém o valor total das compras
func ObterComprasTotal(stateCompra *UserStateCompras, userTokens map[int64]string, chatID int64) (res *ResObterComprasTotal, err error) {
	var (
		resp            *http.Response
		ambiente        = utils.ValidarAmbiente()
		urlObterCompras = ambiente + BaseURLCompras
	)

	urlObterCompras += fmt.Sprintf("/total?data_especifica=%s", *stateCompra.ObterTotalCompras.DataEspecifica)

	if *stateCompra.ObterTotalCompras.Pago {
		urlObterCompras += fmt.Sprintf("&pago=%v", *stateCompra.ObterTotalCompras.Pago)
	}

	if *stateCompra.ObterTotalCompras.UltimaParcela {
		urlObterCompras += fmt.Sprintf("&ultima_parcela=%v", *stateCompra.ObterTotalCompras.UltimaParcela)
	}

	token, ok := userTokens[chatID]
	if !ok {
		return res, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, urlObterCompras, nil)
	if err != nil {
		return res, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		return res, err
	}
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
