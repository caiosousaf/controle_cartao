package cartao

import (
	"github.com/google/uuid"
	"time"
)

// UserStateCartao Struct para armazenar o estado da conversa do usuário para ações de cartões
type UserStateCartao struct {
	CurrentStep     string
	CurrentStepBool bool
	NovoCartaoData  NovoCartao
}

// NovoCartao Struct para armazenar os dados de um novo cartão
type NovoCartao struct {
	ID   string
	Nome string
}

const BaseURLCartoes = "http://localhost:8080/cadastros/cartoes"

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
