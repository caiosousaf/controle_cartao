package postgres

import (
	model "controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
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
