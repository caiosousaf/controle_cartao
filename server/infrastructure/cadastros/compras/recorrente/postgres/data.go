package postgres

import (
	model "controle_cartao/infrastructure/cadastros/compras/recorrente"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBRecorrentes é uma estrutura para acesso aos metodos do banco postgres para
// manipulação dos dados de compras recorrentes
type DBRecorrentes struct {
	DB *sql.DB
}

// ListarComprasRecorrentes lista de forma paginada os dados das compras recorrentes
func (pg *DBRecorrentes) ListarComprasRecorrentes(params *utils.Parametros, usuarioID *uuid.UUID) (res *model.RecorrentesPag, err error) {
	var t model.Recorrentes

	res = new(model.RecorrentesPag)

	campos, _, err := params.ValidFields(&t)
	if err != nil {
		return
	}

	consulta := sq.StatementBuilder.RunWith(pg.DB).
		Select(campos...).
		From("public.t_compras_recorrente TCR").
		Where(sq.Eq{"usuario_id": usuarioID})

	consultaComFiltro := params.CriarFiltros(consulta, map[string]utils.Filtro{
		"ativo": utils.CriarFiltros("TCR.ativo = ?", utils.FlagFiltroEq),
	}).
		PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(params, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Recorrentes), prox, total

	return
}
