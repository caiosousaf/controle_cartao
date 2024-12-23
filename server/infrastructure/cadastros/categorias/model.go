package categorias

import (
	"time"

	"github.com/google/uuid"
)

// Categorias é uma estrutura para definição de modelo de categoria para uso na camada de dados
type Categorias struct {
	ID              *uuid.UUID `sql:"id" apelido:"id"`
	Nome            *string    `sql:"nome" apelido:"nome"`
	UsuarioID       *uuid.UUID `sql:"usuario_id" apelido:"usuario_id"`
	DataCriacao     *time.Time `sql:"data_criacao" apelido:"data_criacao"`
	DataDesativacao *time.Time `sql:"data_desativacao" apelido:"data_desativacao"`
}

// CategoriasPag é uma estrutura para retorno de lista de dados paginada
type CategoriasPag struct {
	Dados []Categorias
	Prox  *bool
	Total *int64
}
