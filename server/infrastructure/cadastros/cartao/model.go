package cartao

import (
	"github.com/google/uuid"
	"time"
)

// Cartao estrutura para definição de modelo de cartão para uso na camada de dados
type Cartao struct {
	ID              *uuid.UUID `sql:"id" apelido:"id"`
	Nome            *string    `sql:"nome" apelido:"nome"`
	DataCriacao     *time.Time `sql:"data_criacao" apelido:"data_criacao"`
	DataDesativacao *time.Time `sql:"data_desativacao" apelido:"data_desativacao"`
	UsuarioID       *uuid.UUID `sql:"usuario_id" apelido:"usuario_id"`
}

// CartaoPag estrutura para retorno de lista de dados paginada
type CartaoPag struct {
	Dados []Cartao
	Prox  *bool
	Total *int64
}
