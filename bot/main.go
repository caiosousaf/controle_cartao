package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Res modela uma resposta para listagem e busca de cartões
type Res struct {
	ID              *uuid.UUID `json:"id" apelido:"id"`
	Nome            *string    `json:"nome" apelido:"nome"`
	DataCriacao     *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataDesativacao *time.Time `json:"data_desativacao" apelido:"data_desativacao"`
}

// ResPag modela uma lista de respostas com suporte para paginação de cartões na listagem
type ResPag struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}

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

	var cartoes ResPag
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

func getCartoes(url string) (cartoes ResPag) {
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

// Res modela uma resposta para listagem e busca de faturas de um cartão
type ResFaturas struct {
	ID             *uuid.UUID `json:"id" apelido:"id"`
	Nome           *string    `json:"nome" apelido:"nome"`
	FaturaCartaoID *uuid.UUID `json:"fatura_cartao_id" apelido:"cartao_id"`
	NomeCartao     *string    `json:"nome_cartao" apelido:"nome_cartao"`
	Status         *string    `json:"status" apelido:"status"`
	DataCriacao    *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataVencimento *string    `json:"data_vencimento" apelido:"data_vencimento"`
}

// ResPag modela uma lista de respostas com suporte para paginação de faturas de cartão na listagem
type ResPagFaturas struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}

func getFaturas(url string) (res ResPagFaturas) {
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

func main() {
	const (
		baseURLCartoes = "http://localhost:8080/cadastros/cartoes"
		baseURLFaturas = "http://localhost:8080/cadastros/cartao/"
	)

	var (
		faturaCartao   = false
		optionsCartoes []string
	)
	// Inicialize o token do seu bot aqui
	token := "MEU_TOKEN"

	// Cria um novo bot com o token fornecido
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Configuração de debug para receber informações detalhadas
	bot.Debug = true

	log.Printf("Autorizado como %s", bot.Self.UserName)

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
		if update.Message == nil { // Ignora atualizações sem mensagem
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		tgbotapi.NewKeyboardButton("ddafsd")

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
		}

		if update.Message.Text == "/cartoes" || update.Message.Text == "cartoes" {
			msgGet := realizarGetString(baseURLCartoes)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgGet)
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		}

		if update.Message.Text == "faturas" || update.Message.Text == "/faturas" {
			cartoes := getCartoes(baseURLCartoes)

			var options []string
			for _, cartao := range cartoes.Dados {
				options = append(options, *cartao.Nome)
			}

			optionsCartoes = options

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Selecione o cartão:")
			msg.ReplyMarkup = keyboard

			// Enviando a mensagem
			_, err := bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}

			faturaCartao = true
		}

		if faturaCartao == true && len(optionsCartoes) != 0 {
			var (
				valueOption *string
				options     []string
			)

			for _, option := range optionsCartoes {
				if update.Message.Text == option {
					valueOption = &option
					break
				}
			}

			if valueOption != nil {
				cartoes := getCartoes(fmt.Sprintf(baseURLCartoes+"?nome_exato=%v", *valueOption))

				faturas := getFaturas(fmt.Sprintf(baseURLFaturas+"%v/faturas", cartoes.Dados[0].ID))

				for _, fatura := range faturas.Dados {
					options = append(options, *fatura.Nome)
				}

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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Selecione a fatura:")
				msg.ReplyMarkup = keyboard

				// Enviando a mensagem
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
			}
		}

	}
}
