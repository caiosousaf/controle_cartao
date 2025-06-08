package usuarios

import (
	"github.com/google/uuid"
	"time"
)

// ReqUsuario modela uma estrutura para a criação de um usuário
type ReqUsuario struct {
	ID    *uuid.UUID `json:"-" apelido:"id"`
	Nome  *string    `json:"nome" apelido:"nome"`
	Email *string    `json:"email" apelido:"email"`
	Senha *string    `json:"senha" apelido:"senha"`
}

// ResCadastroUsuario modela estrutura de resposta em caso de sucesso do cadastro de usuario
type ResCadastroUsuario struct {
	ID    *uuid.UUID `json:"id"`
	Token string     `json:"token"`
}

// Res modela estrutura de resposta para sucesso de login
type Res struct {
	Token string `json:"token"`
}

// ReqUsuarioLogin modela uma estrutura para o login de um usuário
type ReqUsuarioLogin struct {
	Email *string `json:"email" apelido:"email" binding:"required"`
	Senha *string `json:"senha" apelido:"senha" binding:"required"`
}

// ReqAlterarSenhaUsuario modela uma estrutura para alterar a senha de um usuário
type ReqAlterarSenhaUsuario struct {
	Email      *string `json:"email" binding:"required"`
	EmailNovo  *string `json:"email_novo" binding:"required"`
	SenhaAtual *string `json:"senha_atual" binding:"required"`
	SenhaNova  *string `json:"senha_nova" binding:"required"`
}

// ResUsuario modela uma estrutura de resposta com os dados de usuário
type ResUsuario struct {
	Nome        *string    `json:"nome" apelido:"nome"`
	Email       *string    `json:"email" apelido:"email"`
	DataCriacao *time.Time `json:"data_criacao" apelido:"data_criacao"`
}
