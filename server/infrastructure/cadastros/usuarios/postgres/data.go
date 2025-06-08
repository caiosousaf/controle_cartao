package postgres

import (
	model "controle_cartao/infrastructure/cadastros/usuarios"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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

// BuscarUsuarioLogin é responsável por buscar um usuário no banco de dados
func (pg *DBUsuario) BuscarUsuarioLogin(email *string) (res *model.Usuario, err error) {
	res = new(model.Usuario)

	if err = sq.StatementBuilder.RunWith(pg.DB).Select("id", "nome", "email", "senha", "data_criacao", "data_desativacao").
		From("public.t_usuarios").
		Where("email = $1", email).
		Where("data_desativacao ISNULL").
		Scan(&res.ID, &res.Nome, &res.Email, &res.Senha, &res.DataCriacao, &res.DataDesativacao); err != nil {
		return nil, err
	}

	return
}

// BuscarUsuario é responsável por buscar um usuário no banco de dados a partir do seu ID
func (pg *DBUsuario) BuscarUsuario(usuarioID *uuid.UUID) (res *model.Usuario, err error) {
	res = new(model.Usuario)

	if err = sq.StatementBuilder.RunWith(pg.DB).Select("id", "nome", "email", "senha", "data_criacao", "data_desativacao").
		From("public.t_usuarios").
		Where("id = $1", usuarioID).
		Where("data_desativacao ISNULL").
		Scan(&res.ID, &res.Nome, &res.Email, &res.Senha, &res.DataCriacao, &res.DataDesativacao); err != nil {
		return nil, err
	}

	return
}

// AtualizarSenhaUsuario é responsável por atualizar a senha do usuário
func (pg *DBUsuario) AtualizarSenhaUsuario(novaSenha, email *string, usuarioID *uuid.UUID) error {
	if _, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_usuarios").
		Set("senha", novaSenha).
		Set("email", email).
		Where(sq.Eq{
			"id": usuarioID,
		}).
		PlaceholderFormat(sq.Dollar).
		Exec(); err != nil {
		return err
	}

	return nil
}
