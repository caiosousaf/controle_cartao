package cartao

import (
	"github.com/google/uuid"
	"time"
)

// Req modela uma requisição para a criação de um cartão
type Req struct {
	Nome      *string    `json:"nome" apelido:"nome"`
	UsuarioID *uuid.UUID `json:"-" apelido:"usuario_id"`
}

// ReqAtualizar modela uma requisição para a atualização de um cartão
type ReqAtualizar struct {
	Nome *string `json:"nome" apelido:"nome"`
}

// Res modela uma resposta para listagem e busca de cartões
type Res struct {
	ID              *uuid.UUID `json:"id" apelido:"id"`
	Nome            *string    `json:"nome" apelido:"nome"`
	UsuarioID       *uuid.UUID `json:"usuario_id" apelido:"usuario_id"`
	DataCriacao     *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataDesativacao *time.Time `json:"data_desativacao" apelido:"data_desativacao"`
}

// ResPag modela uma lista de respostas com suporte para paginação de cartões na listagem
type ResPag struct {
	Dados []Res  `json:"dados,omitempty"`
	Prox  *bool  `json:"prox,omitempty"`
	Total *int64 `json:"total,omitempty"`
}
