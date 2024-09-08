package postgres

import (
	model "controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBCategoria é uma estrutura para acesso aos métodos do banco postgres para manipulação das categorias
type DBCategoria struct {
	DB *sql.DB
}

// CadastrarCategoria cadastra uma nova categoria no banco de dados
func (pg *DBCategoria) CadastrarCategoria(req *model.Categorias) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_categoria_compra").
		Columns("nome", "usuario_id").
		Values(req.Nome, req.UsuarioID).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).
		Scan(&req.ID); err != nil {
		return
	}

	return
}

// AtualizarCategoria atualiza uma categoria no banco de dados
func (pg *DBCategoria) AtualizarCategoria(req *model.Categorias, idCategoria *uuid.UUID) (err error) {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_categoria_compra").
		Set("nome", req.Nome).
		Where(sq.Eq{"id": idCategoria,
			"data_desativacao": nil,
			"usuario_id":       req.UsuarioID,
		}).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("Categoria não foi encontrada, ou não existe no banco de dados")
	}

	return err
}

// ListarCategorias lista as categorias cadastradas
func (pg *DBCategoria) ListarCategorias(params *utils.Parametros) (res *model.CategoriasPag, err error) {
	var t model.Categorias

	res = new(model.CategoriasPag)

	campos, _, err := params.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_categoria_compra TCC")

	if !params.TemFiltro("ativo") {
		consultaSql = consultaSql.Where(sq.Eq{"TCC.data_desativacao": nil})
	}

	consultaComFiltro := params.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"categoria_id": utils.CriarFiltros("TCC.id = ?::UUID", utils.FlagFiltroEq),
		"nome_exato":   utils.CriarFiltros("TCC.nome = ?::VARCHAR", utils.FlagFiltroEq),
		"ativo":        utils.CriarFiltros("TCC.data_desativacao IS NULL = ?::BOOLEAN", utils.FlagFiltroEq),
		"usuario_id":   utils.CriarFiltros("TCC.usuario_id = ?::UUID", utils.FlagFiltroEq),
	}).
		PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(params, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Categorias), prox, total

	return
}

// RemoverCategoria remove uma categoria
func (pg *DBCategoria) RemoverCategoria(idCategoria, usuarioID *uuid.UUID) error {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_categoria_compra").
		Set("data_desativacao", "NOW()").
		Where(sq.Eq{
			"id":         idCategoria,
			"usuario_id": usuarioID,
		}).PlaceholderFormat(sq.Dollar).Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("Categoria não foi encontrada, já se encontra desativada ou não existe")
	}

	return err
}

// ReativarCategoria reativa uma categoria
func (pg *DBCategoria) ReativarCategoria(idCategoria, usuarioID *uuid.UUID) error {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_categoria_compra").
		Set("data_desativacao", nil).
		Where(sq.Eq{
			"id":         idCategoria,
			"usuario_id": usuarioID,
		},
			sq.NotEq{
				"data_desativacao": nil,
			}).PlaceholderFormat(sq.Dollar).Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("Categoria não foi encontrada, já se encontra ativada ou não existe no banco de dados")
	}

	return err
}

// BuscarCategoria busca uma categoria de acordo com o id dela
func (pg *DBCategoria) BuscarCategoria(idCategoria, usuarioID *uuid.UUID) (res *model.Categorias, err error) {
	res = new(model.Categorias)

	if err = sq.StatementBuilder.RunWith(pg.DB).Select(`TCC.id, TCC.nome, TCC.data_criacao, TCC.data_desativacao, TCC.usuario_id`).
		From("public.t_categoria_compra TCC").
		Where(sq.Eq{
			"id":               idCategoria,
			"usuario_id":       usuarioID,
			"data_desativacao": nil,
		}).
		PlaceholderFormat(sq.Dollar).
		Scan(&res.ID, &res.Nome, &res.DataCriacao, &res.DataDesativacao, &res.UsuarioID); err != nil {
		return res, err
	}
	return
}
