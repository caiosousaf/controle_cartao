package usuarios

import (
	"controle_cartao/infrastructure/cadastros/usuarios"
	"controle_cartao/infrastructure/cadastros/usuarios/postgres"
	"database/sql"
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

// BuscarUsuario é um gerenciador de fluxo de dados para buscar um usuário no banco de dados
func (r *repo) BuscarUsuario(email *string) (*usuarios.Usuario, error) {
	return r.Data.BuscarUsuario(email)
}
