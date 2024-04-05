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
func (pg *DBCategoria) RemoverCategoria(idCategoria *uuid.UUID) error {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_categoria_compra").
		Set("data_desativacao", "NOW()").
		Where(sq.Eq{
			"id": idCategoria,
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
func (pg *DBCategoria) ReativarCategoria(idCategoria *uuid.UUID) error {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_categoria_compra").
		Set("data_desativacao", nil).
		Where(sq.Eq{
			"id": idCategoria,
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
