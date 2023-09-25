package postgres

import (
	model "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBFatura é uma estrutura base para acesso aos metodos do banco postgres para manipulação das faturas
type DBFatura struct {
	DB *sql.DB
}

// ListarFaturas lista de forma paginada os dados das faturas de um cartão no banco de dados
func (pg *DBFatura) ListarFaturasCartao(p *utils.Parametros, id *uuid.UUID) (res *model.FaturaPag, err error) {
	var t model.Fatura

	res = new(model.FaturaPag)

	campos, _, err := p.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC on TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{"TFC.fatura_cartao_id": id})

	consultaComFiltro := p.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"nome_exato": utils.CriarFiltros("LOWER(TF.nome) = LOWER(?)", utils.FlagFiltroEq),
	}).PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(p, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Fatura), prox, total

	return
}
