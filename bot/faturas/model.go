package faturas

import (
	"bot_controle_cartao/utils"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// UserState Struct para armazenar o estado da conversa do usuário
type UserState struct {
	ChatID          int64
	CurrentStep     string
	CurrentStepBool bool
	NewInvoiceData  NewInvoice
}

// NewInvoice Struct para armazenar os dados de uma nova fatura
type NewInvoice struct {
	Title   string
	Amount  float64
	DueDate string
}

type UserStepComprasFatura struct {
	ComprasFatura bool
	Cartoes       []string
	Opcao         *string
	Fatura        Fatura
}

type Fatura struct {
	Position *int
	ID       *uuid.UUID
	Nome     *string
	CartaoID *uuid.UUID
}

type ReqAtualizarStatus struct {
	Status *string `json:"status"`
}

var url = utils.ValidarSistema()

var (
	BaseURLFaturas = fmt.Sprintf("%s/cadastros/cartao/", url)
	BaseURLFatura  = fmt.Sprintf("%s/cadastros/fatura/", url)
)

// Res modela uma resposta para listagem e busca de faturas de um cartão
type Res struct {
	ID             *uuid.UUID `json:"id" apelido:"id"`
	Nome           *string    `json:"nome" apelido:"nome"`
	FaturaCartaoID *uuid.UUID `json:"fatura_cartao_id" apelido:"cartao_id"`
	NomeCartao     *string    `json:"nome_cartao" apelido:"nome_cartao"`
	Status         *string    `json:"status" apelido:"status"`
	DataCriacao    *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataVencimento *string    `json:"data_vencimento" apelido:"data_vencimento"`
}

// ResPagFaturas modela uma lista de respostas com suporte para paginação de faturas de cartão na listagem
type ResPagFaturas struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}
