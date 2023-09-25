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

// ListarFaturasCartao lista de forma paginada os dados das faturas de um cartão no banco de dados
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
		Where(sq.Eq{
			"TFC.fatura_cartao_id": id,
			"TC.data_desativacao":  nil,
		})

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

// BuscarFaturaCartao busca os dados de uma fatura de um cartão no banco de dados dado os ID's fornecidos
func (pg *DBFatura) BuscarFaturaCartao(idFatura, idCartao *uuid.UUID) (res *model.Fatura, err error) {
	res = new(model.Fatura)
	if err = sq.StatementBuilder.RunWith(pg.DB).Select(`TFC.id, TFC.nome, TFC.fatura_cartao_id,
				TC.nome, TFC.data_criacao, TFC.data_vencimento`).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC on TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.fatura_cartao_id": idCartao,
			"TFC.id":               idFatura,
			"TC.data_desativacao":  nil,
		}).PlaceholderFormat(sq.Dollar).
		Scan(&res.ID, &res.Nome, &res.FaturaCartaoID, &res.NomeCartao, &res.DataCriacao, &res.DataVencimento); err != nil {
		return res, err
	}

	return
}
