package usuarios

import "github.com/google/uuid"

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
