package faturas

import (
	"github.com/google/uuid"
	"time"
)

// Req modela uma requisição para a criação de uma fatura
type Req struct {
	Nome           *string    `json:"nome" apelido:"nome"`
	FaturaCartaoID *uuid.UUID `json:"cartao_id" apelido:"cartao_id"`
	DataVencimento *string    `json:"data_vencimento" apelido:"data_vencimento"`
}

// ReqAtualizar modela uma requisição para a atualização de uma fatura
type ReqAtualizar struct {
	Nome           *string `json:"nome" apelido:"nome"`
	DataVencimento *string `json:"data_vencimento" apelido:"data_vencimento"`
}

// Res modela uma resposta para listagem e busca de faturas de um cartão
type Res struct {
	ID             *uuid.UUID `json:"id" apelido:"id"`
	Nome           *string    `json:"nome" apelido:"nome"`
	FaturaCartaoID *uuid.UUID `json:"fatura_cartao_id" apelido:"cartao_id"`
	NomeCartao     *string    `json:"nome_cartao" apelido:"nome_cartao"`
	DataCriacao    *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataVencimento *string    `json:"data_vencimento" apelido:"data_vencimento"`
}

// ResPag modela uma lista de respostas com suporte para paginação de faturas de cartão na listagem
type ResPag struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}
