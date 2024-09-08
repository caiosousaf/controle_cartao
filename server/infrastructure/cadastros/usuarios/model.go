package usuarios

import (
	"github.com/google/uuid"
	"time"
)

// Usuario é a estruturea para definição de modelo de usuário
type Usuario struct {
	ID              *uuid.UUID `sql:"id" apelido:"id"`
	Nome            *string    `sql:"nome" apelido:"nome"`
	Email           *string    `sql:"email" apelido:"email"`
	Senha           *string    `sql:"senha" apelido:"senha"`
	DataCriacao     *time.Time `sql:"data_criacao"`
	DataDesativacao *time.Time `sql:"data_desativacao"`
}
