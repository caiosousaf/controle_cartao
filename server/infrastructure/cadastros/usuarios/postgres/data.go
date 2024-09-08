package postgres

import (
	model "controle_cartao/infrastructure/cadastros/usuarios"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

// DBUsuario é uma estrutura base para acesso aos metodos do banco postgresql
type DBUsuario struct {
	DB *sql.DB
}

// CadastrarUsuario é responsável por cadastrar um novo usuario no banco de dados
func (pg *DBUsuario) CadastrarUsuario(req *model.Usuario) error {
	if err := sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_usuarios").
		Columns("nome", "email", "senha").
		Values(req.Nome, req.Email, req.Senha).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).
		Scan(&req.ID); err != nil {
		return err
	}

	return nil
}

// BuscarUsuario é responsável por buscar um usuário no banco de dados
func (pg *DBUsuario) BuscarUsuario(email *string) (res *model.Usuario, err error) {
	res = new(model.Usuario)

	if err = sq.StatementBuilder.RunWith(pg.DB).Select("id", "nome", "email", "senha", "data_criacao", "data_desativacao").
		From("public.t_usuarios").
		Where("email = $1", email).
		Scan(&res.ID, &res.Nome, &res.Email, &res.Senha, &res.DataCriacao, &res.DataDesativacao); err != nil {
		return nil, err
	}

	return
}
