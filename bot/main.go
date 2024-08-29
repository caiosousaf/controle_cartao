package main

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/categorias"
	"bot_controle_cartao/compras"
	"bot_controle_cartao/faturas"
	"bot_controle_cartao/utils"
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

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
		Fatura:  faturas.Fatura{},
	}
	userCompras := &compras.UserStateCompras{
		FaturaID:          nil,
		CurrentStep:       nil,
		CurrentStepBool:   false,
		NovaCompraData:    compras.NovaCompra{},
		ObterTotalCompras: compras.ObterCompras{},
	}
	userStatusFatura := &faturas.ReqAtualizarStatus{Status: nil}

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
						if err.Error() == utils.ErroPdfVazio {
							msgErrArquivoVazio := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Esse cartão não possui compras no ano selecionado")
							utils.EnviaMensagem(bot, msgErrArquivoVazio)
						} else {
							log.Panic(err)
						}
					}
				}
			}

			if AcaoAnterior == "faturas" {
				faturas.ProcessarCasosStepComprasFatura(userCompraFaturas, userStatusFatura, bot, update)
			}

			if AcaoAnterior == "compras" {
				switch *userCompras.CurrentStep {
				case "selecionar_fatura":
					faturasCartao := faturas.ListarFaturas(fmt.Sprintf(faturas.BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data))

					faturas.EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, userCompras, update.CallbackQuery)
				case "fatura_selecionada":
					categoriasCompras := categorias.ListarCategorias(categorias.BaseURLCategoria)

					compras.EnviarOpcoesCategoriasCompras(bot, update.CallbackQuery.Message.Chat.ID, &categoriasCompras, userCompras, update.CallbackQuery)
				case "categoria_selecionada":
					compras.InicioCriacaoCompra(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userCompras)

					AcaoAnterior = "cadastro_de_compra"
				}
			}
		} else if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				faturas.EnviaMensagemBoasVindas(bot, update.Message.Chat.ID)

				// Criando um teclado de resposta
				buttonOpcao1 := tgbotapi.NewKeyboardButton("Cartões")
				buttonOpcao2 := tgbotapi.NewKeyboardButton("Faturas")
				buttonOpcao3 := tgbotapi.NewKeyboardButton("Compras")
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

			if update.Message.Text == "/Cartões" || update.Message.Text == "Cartões" || AcaoAnterior == "cartoes" {
				cartao.ProcessoAcoesCartoes(bot, update.Message, userStatesCartao)

				AcaoAnterior = "cartoes"
			}

			if update.Message.Text == "Faturas" || update.Message.Text == "/Faturas" || AcaoAnterior == "faturas" {
				faturas.ProcessoAcoesFaturas(bot, update.Message, userStates, userCompraFaturas)

				AcaoAnterior = "faturas"
			}

			if update.Message.Text == "Compras" || update.Message.Text == "/Compras" || AcaoAnterior == "compras" {
				compras.ProcessoAcoesCompras(bot, update.Message, userCompras, AcaoAnterior)

				AcaoAnterior = "compras"
			}

			if AcaoAnterior == "cadastro_de_compra" {
				if userCompras.CurrentStep != nil {
					compras.ProcessoAcoesCadastroCompra(bot, update.Message, userCompras)
				}
			}
		}
	}

	go func() {
		log.Println("Starting HTTP server on port 8079...")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Bot is running"))
		})
		if err := http.ListenAndServe(":8079", nil); err != nil {
			log.Fatal(err)
		}
	}()
}
