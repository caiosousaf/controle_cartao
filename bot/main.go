package main

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/compras"
	"bot_controle_cartao/faturas"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func init() {
	if err := godotenv.Load("server/.env"); err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}
}

func main() {

	// Inicialize o token do seu bot aqui
	token := os.Getenv("TOKEN")

	// Cria um novo bot com o token fornecido
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Configuração de debug para receber informações detalhadas
	bot.Debug = true

	log.Printf("Autorizado como %s", bot.Self.UserName)

	userStates := make(map[int64]*faturas.UserState)
	userStatesCartao := &cartao.UserStateCartao{
		CurrentStep:     "",
		CurrentStepBool: false,
		NovoCartaoData:  cartao.NovoCartao{},
	}
	userCompraFaturas := &faturas.UserStepComprasFatura{
		Cartoes: []string{},
		Opcao:   nil,
	}

	var (
		AcaoAnterior string
	)

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

			if AcaoAnterior == "cartoes" {
				switch userStatesCartao.CurrentStep {
				case "selecionar_ano":
					userStatesCartao.NovoCartaoData.ID = update.CallbackQuery.Data

					cartao.EnviarOpcoesAno(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userStatesCartao)
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

			if AcaoAnterior == "faturas" {
				switch *userCompraFaturas.Opcao {
				case "fatura_selecionada":
					faturasCartao := faturas.ListarFaturas(fmt.Sprintf(faturas.BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data))

					faturas.EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, update.CallbackQuery)
				case "cartao_fatura_selecionado":
					faturas.ProcessCallbackQuery(bot, update.CallbackQuery)
				}
			}
		} else if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				faturas.EnviaMensagemBoasVindas(bot, update.Message.Chat.ID)

				// Criando um teclado de resposta
				buttonOpcao1 := tgbotapi.NewKeyboardButton("cartoes")
				buttonOpcao2 := tgbotapi.NewKeyboardButton("faturas")
				buttonOpcao3 := tgbotapi.NewKeyboardButton("compras")
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
				AcaoAnterior = "start"
			}

			if update.Message.Text == "/cartoes" || update.Message.Text == "cartoes" || AcaoAnterior == "cartoes" {
				cartao.ProcessoAcoesCartoes(bot, update.Message, userStatesCartao)

				AcaoAnterior = "cartoes"
			}

			if update.Message.Text == "faturas" || update.Message.Text == "/faturas" || AcaoAnterior == "faturas" {
				faturas.ProcessoAcoesFaturas(bot, update.Message, userStates, userCompraFaturas)

				AcaoAnterior = "faturas"
			}
		}
	}
}
