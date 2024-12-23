package postgres

import (
	model "controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBCartao é uma estrutura base para acesso aos metodos do banco postgres para manipulação de cartões
type DBCartao struct {
	DB *sql.DB
}

// CadastrarCartao cadastra um novo cartao no banco de dados
func (pg *DBCartao) CadastrarCartao(req *model.Cartao) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_cartao").
		Columns("nome", "usuario_id").
		Values(*req.Nome, req.UsuarioID).
		PlaceholderFormat(sq.Dollar).
		Suffix(`RETURNING "id"`).
		Scan(&req.ID); err != nil {
		return err
	}

	return
}

// ListarCartoes lista de forma paginada os dados dos cartões cadastrados no banco de dados
func (pg *DBCartao) ListarCartoes(p *utils.Parametros) (res *model.CartaoPag, err error) {
	var t model.Cartao

	res = new(model.CartaoPag)

	campos, _, err := p.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_cartao TC").
		Where("TC.data_desativacao IS NULL")

	consultaComFiltro := p.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"nome_exato": utils.CriarFiltros("lower(TC.nome) = lower(?)", utils.FlagFiltroEq),
		"ativo":      utils.CriarFiltros("data_desativacao IS NULL = ?", utils.FlagFiltroEq),
		"usuario_id": utils.CriarFiltros("usuario_id = ?::UUID", utils.FlagFiltroEq),
	}).PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(p, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Cartao), prox, total

	return
}

// BuscarCartao busca os dados de um cartão no banco de dados dado o seu id
func (pg *DBCartao) BuscarCartao(id, usuarioID *uuid.UUID) (res *model.Cartao, err error) {
	res = new(model.Cartao)
	if err = sq.StatementBuilder.RunWith(pg.DB).Select("id, nome, usuario_id, data_criacao, data_desativacao").
		From("public.t_cartao").
		Where("id = $1", id).
		Where("usuario_id = $2", usuarioID).
		Where("data_desativacao ISNULL").Scan(&res.ID, &res.Nome, &res.UsuarioID, &res.DataCriacao, &res.DataDesativacao); err != nil {
		return res, err
	}

	return
}

// AtualizarCartao atualiza os dados de um cartão no banco de dados dado o seu id
func (pg *DBCartao) AtualizarCartao(req *model.Cartao, id, usuarioID *uuid.UUID) (err error) {
	updateBuilder := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_cartao").
		Set("nome", req.Nome).
		Where(sq.Eq{
			"id":               id,
			"usuario_id":       usuarioID,
			"data_desativacao": nil,
		}).
		PlaceholderFormat(sq.Dollar)

	result, err := updateBuilder.Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("ID do cartão inexistente")
	}

	return
}

// RemoverCartao desativa os dados de um cartão no banco de dados dado o seu id
func (pg *DBCartao) RemoverCartao(id, usuarioID *uuid.UUID) (err error) {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_cartao").
		Set("data_desativacao", "NOW()").
		Where(sq.Eq{
			"id":         id,
			"usuario_id": usuarioID,
		}).
		PlaceholderFormat(sq.Dollar).
		Exec()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("ID do cartão inexistente")
	}

	return
}

// ReativarCartao reativa os dados de um cartão no banco de dados dado o seu id
func (pg *DBCartao) ReativarCartao(id, usuarioID *uuid.UUID) (err error) {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_cartao").
		Set("data_desativacao", nil).
		Where(sq.Eq{
			"id":         id,
			"usuario_id": usuarioID,
		}).
		PlaceholderFormat(sq.Dollar).
		Exec()

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("ID do cartão inexistente")
	}

	return
}
