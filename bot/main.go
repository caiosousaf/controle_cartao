package main

import (
	"bot_controle_cartao/cartao"
	"bot_controle_cartao/categorias"
	"bot_controle_cartao/compras"
	"bot_controle_cartao/faturas"
	"bot_controle_cartao/usuarios"
	"bot_controle_cartao/utils"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
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
	userStatesCartao := make(map[int64]*cartao.UserStateCartao)
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
		userTokens   = make(map[int64]string)
		email        *string
	)

	// Configuração de atualização com o webhook ou polling
	// Aqui, estamos usando a opção de polling para obter atualizações
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
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

	// Loop pelas atualizações recebidas do bot
	for update := range updates {
		if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Message.Text)

			if AcaoAnterior == "cartoes" {
				switch userStatesCartao[update.CallbackQuery.Message.Chat.ID].CurrentStep {
				case "selecionar_ano":
					userStatesCartao[update.CallbackQuery.Message.Chat.ID].NovoCartaoData.ID = update.CallbackQuery.Data

					cartao.EnviarOpcoesAno(bot, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery, userStatesCartao)
				case "ano_selecionado":
					idCartaoUUID, err := uuid.Parse(userStatesCartao[update.CallbackQuery.Message.Chat.ID].NovoCartaoData.ID)
					if err != nil {
						log.Panic(err)
					}

					edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("Cartão Selecionado: %s", update.CallbackQuery.Data))
					edit.ReplyMarkup = nil

					_, err = bot.Send(edit)
					if err != nil {
						log.Panic(err)
					}

					pdfContent, err := compras.ObterComprasPdf(nil, &idCartaoUUID, userTokens, update.CallbackQuery.Message.Chat.ID)
					if err != nil {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
						utils.EnviaMensagem(bot, msg)
						return
					}

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
				faturas.ProcessarCasosStepComprasFatura(userCompraFaturas, userStatusFatura, bot, update, userTokens)
			}

			if AcaoAnterior == "compras" {
				switch *userCompras.CurrentStep {
				case "selecionar_fatura":
					faturasCartao, err := faturas.ListarFaturas(fmt.Sprintf(faturas.BaseURLFaturas+"%s/faturas", update.CallbackQuery.Data), userTokens, update.CallbackQuery.Message.Chat.ID)
					if err != nil {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
						utils.EnviaMensagem(bot, msg)
						return
					}

					faturas.EnviarOpcoesFaturas(bot, update.CallbackQuery.Message.Chat.ID, &faturasCartao, userCompraFaturas, userCompras, update.CallbackQuery, userTokens)
				case "fatura_selecionada":
					categoriasCompras, err := categorias.ListarCategorias(categorias.BaseURLCategoria, update.CallbackQuery.Message.Chat.ID, userTokens)
					if err != nil {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
						utils.EnviaMensagem(bot, msg)
						return
					}

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
				buttonOpcao4 := tgbotapi.NewKeyboardButton("Usuarios")

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
				cartao.ProcessoAcoesCartoes(bot, update.Message, userStatesCartao, userTokens)

				AcaoAnterior = "cartoes"
			}

			if update.Message.Text == "Faturas" || update.Message.Text == "/Faturas" || AcaoAnterior == "faturas" {
				faturas.ProcessoAcoesFaturas(bot, update.Message, userStates, userCompraFaturas, userTokens)

				AcaoAnterior = "faturas"
			}

			if update.Message.Text == "Compras" || update.Message.Text == "/Compras" || AcaoAnterior == "compras" {
				compras.ProcessoAcoesCompras(bot, update.Message, userCompras, AcaoAnterior, userTokens)

				AcaoAnterior = "compras"
			}

			if AcaoAnterior == "cadastro_de_compra" {
				if userCompras.CurrentStep != nil {
					compras.ProcessoAcoesCadastroCompra(bot, update.Message, userCompras, userTokens)
				}
			}

			if update.Message.Text == "Usuarios" {
				usuarios.GerarOpcoesAcoesUsuarios(bot, update.Message)
			}

			if AcaoAnterior == "login_username" || AcaoAnterior == "login_password" || update.Message.Text == "Login" {
				if AcaoAnterior == "login_username" && update.Message.Text != "" {
					email = &update.Message.Text
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Agora envie sua senha.")
					bot.Send(msg)

					AcaoAnterior = "login_password"
				} else if AcaoAnterior == "login_password" {
					senha := update.Message.Text

					// Realize o login
					err := usuarios.Login(bot, update.Message.Chat.ID, email, senha, userTokens)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Erro no login: %v", err))
						bot.Send(msg)
					}

					AcaoAnterior = ""
				}

				switch update.Message.Text {
				case "Login":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Por favor, envie seu email de usuário.")
					bot.Send(msg)

					AcaoAnterior = "login_username"
				}
			}
		}
	}

}
