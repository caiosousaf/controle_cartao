package usuarios

import (
	"controle_cartao/infrastructure/cadastros/usuarios"
	"controle_cartao/infrastructure/cadastros/usuarios/postgres"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBUsuario
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{Data: &postgres.DBUsuario{DB: novoDB}}
}

// CadastrarUsuario é um gerenciador de fluxo de dados para cadastrar um novo usuário no banco de dados
func (r *repo) CadastrarUsuario(req *usuarios.Usuario) error {
	return r.Data.CadastrarUsuario(req)
}

// BuscarUsuarioLogin é um gerenciador de fluxo de dados para buscar um usuário no banco de dados
func (r *repo) BuscarUsuarioLogin(email *string) (*usuarios.Usuario, error) {
	return r.Data.BuscarUsuarioLogin(email)
}

// AtualizarSenhaUsuario é um gerenciador de fluxo de dados para atualizar a senha do usuário
func (r *repo) AtualizarSenhaUsuario(novaSenha, email *string, usuarioID *uuid.UUID) error {
	return r.Data.AtualizarSenhaUsuario(novaSenha, email, usuarioID)
}

// BuscarUsuario é um gerenciador de fluxo de dados para buscar um usuário no banco de dados
func (r *repo) BuscarUsuario(usuarioID *uuid.UUID) (*usuarios.Usuario, error) {
	return r.Data.BuscarUsuario(usuarioID)
}
