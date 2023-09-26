package postgres

import (
	model "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
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

// ObterProximasFaturas obtém as próximas possíveis datas de faturas de um cartão no banco de dados
func (pg *DBFatura) ObterProximasFaturas(qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error) {
	var (
		data time.Time
		mes  int
	)
	consultaSelect := fmt.Sprintf(`GENERATE_SERIES(TFC.data_vencimento, TFC.data_vencimento + INTERVAL '%d months' - INTERVAL '1 month',
                		INTERVAL '1 month')::DATE, TC.id,
    EXTRACT(MONTH FROM GENERATE_SERIES(TFC.data_vencimento, TFC.data_vencimento + INTERVAL '%d months' - INTERVAL '1 month',
    INTERVAL '1 month')::DATE) AS numero_mes`, *qtd_parcelas, *qtd_parcelas)

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(consultaSelect).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC ON TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.id": idFatura,
		}).PlaceholderFormat(sq.Dollar)

	rows, err := consultaSql.Query()
	if err != nil {
		return datas, meses, idCartao, err
	}

	for rows.Next() {
		err := rows.Scan(&data, &idCartao, &mes)
		if err != nil {
			return datas, meses, idCartao, err
		}

		// Formata a data e transforma em string
		formattedDate := data.Format("2006-01-02")
		mesFormatado := utils.NumeroParaNomeMes(mes)

		datas = append(datas, formattedDate)
		meses = append(meses, mesFormatado)
	}

	return
}