package categorias

import (
	"github.com/google/uuid"
	"time"
)

// ReqCategoria modela uma estrutura para a criação de uma categoria
type ReqCategoria struct {
	Nome *string `json:"nome" apelido:"nome"`
}

// ResCategorias modela uma resposta para listagem e busca de categorias
type ResCategorias struct {
	ID              *uuid.UUID `json:"id" apelido:"id"`
	Nome            *string    `json:"nome" apelido:"nome"`
	DataCriacao     *time.Time `json:"data_criacao" apelido:"data_criacao"`
	DataDesativacao *time.Time `json:"data_desativacao" apelido:"data_desativacao"`
}

// ResCategoriasPag modela uma lista de respostas com suporte para paginação de categorias
type ResCategoriasPag struct {
	Dados []ResCategorias `json:"dados,omitempty"`
	Prox  *bool           `json:"prox,omitempty"`
	Total *int64          `json:"total,omitempty"`
}
