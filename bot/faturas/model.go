package faturas

import (
	"github.com/google/uuid"
	"time"
)

// Struct para armazenar o estado da conversa do usuário
type UserState struct {
	ChatID          int64
	CurrentStep     string
	CurrentStepBool bool
	NewInvoiceData  NewInvoice
}

// Struct para armazenar os dados de uma nova fatura
type NewInvoice struct {
	Title   string
	Amount  float64
	DueDate string
}

type UserStepComprasFatura struct {
	ComprasFatura bool
	Cartoes       []string
	Opcao         *string
}

var (
	AcaoAnterior string
)

const (
	BaseURLCartoes = "http://localhost:8080/cadastros/cartoes"
	BaseURLFaturas = "http://localhost:8080/cadastros/cartao/"
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

// ResPag modela uma lista de respostas com suporte para paginação de faturas de cartão na listagem
type ResPagFaturas struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}
