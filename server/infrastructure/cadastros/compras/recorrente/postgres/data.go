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

// ObterFaturaCartaoGeral obtém o id da fatura do cartão geral
func (pg *DBRecorrentes) ObterFaturaCartaoGeral(usuarioID *uuid.UUID) (faturaID *uuid.UUID, err error) {
	consulta := sq.StatementBuilder.RunWith(pg.DB).
		Select("TFC.id").
		From("t_fatura_cartao tfc").
		Join("t_cartao TC ON TFC.FATURA_CARTAO_ID = TC.ID").
		Where(sq.Eq{"usuario_id": usuarioID}).
		Where("TO_CHAR(data_vencimento, 'MM/YYYY') = TO_CHAR(NOW()::DATE, 'MM/YYYY')").
		Where("TC.nome = 'Cartão Geral'").
		Where(sq.NotEq{"status": "Pago"}).
		Where(`tfc.id NOT IN (
			SELECT tcf.compra_fatura_id
			FROM t_compras_fatura tcf
			WHERE tcf.compra_fatura_id = tfc.id
			AND tcf.recorrente )`).PlaceholderFormat(sq.Dollar)

	if err = consulta.Scan(&faturaID); err != nil {
		return
	}

	return
}

// ObterPrevisaoGastos obtém a previsão de gastos dos próximos 3 meses com as compras recorrentes
func (pg *DBRecorrentes) ObterPrevisaoGastos(usuarioID *uuid.UUID) (gastos *model.PrevisaoGastosPag, err error) {
	gastos = new(model.PrevisaoGastosPag)

	query := `
	WITH compras_base AS (
  		SELECT
    		TFC.data_vencimento,
    		TCF.valor_parcela,
    		TCF.recorrente
  		FROM PUBLIC.T_COMPRAS_FATURA TCF
  		JOIN PUBLIC.T_FATURA_CARTAO TFC ON TFC.id = TCF.compra_fatura_id
  		JOIN PUBLIC.T_CARTAO TC ON TC.id = TFC.fatura_cartao_id
  		WHERE TC.usuario_id = $1
),
recorrentes_expandido AS (
  SELECT data_vencimento, valor_parcela FROM compras_base WHERE recorrente = FALSE
  UNION ALL
  SELECT data_vencimento, valor_parcela FROM compras_base WHERE recorrente = TRUE
  UNION ALL
  SELECT data_vencimento + INTERVAL '1 month', valor_parcela FROM compras_base WHERE recorrente = TRUE
  UNION ALL
  SELECT DATA_VENCIMENTO + INTERVAL '2 month', valor_parcela FROM compras_base WHERE recorrente = TRUE
)
	SELECT
  		TO_CHAR(data_vencimento, 'MM/YYYY') AS mes_ano,
  		CAST(SUM(valor_parcela) AS NUMERIC(10,2)) AS valor_total
	FROM RECORRENTES_EXPANDIDO
	WHERE TO_CHAR(data_vencimento, 'MM/YYYY') IN (
  		TO_CHAR(NOW(), 'MM/YYYY'),
  		TO_CHAR(NOW() + INTERVAL '1 month', 'MM/YYYY'),
  		TO_CHAR(NOW() + INTERVAL '2 month', 'MM/YYYY')
)
	GROUP BY TO_CHAR(data_vencimento, 'MM/YYYY')
	ORDER BY MIN(data_vencimento)
`

	rows, err := pg.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gastos.Dados = make([]model.PrevisaoGastos, 0)
	for rows.Next() {
		var p model.PrevisaoGastos
		if err := rows.Scan(&p.MesAno, &p.Valor); err != nil {
			return nil, err
		}
		gastos.Dados = append(gastos.Dados, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

// CadastrarCompraRecorrente cadastra as compras recorrentes
func (pg *DBRecorrentes) CadastrarCompraRecorrente(req *model.ComprasRecorrentes) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_compras_fatura").
		Columns("nome", "descricao", "local_compra", "compra_categoria_id", "valor_parcela", "parcela_atual", "qtd_parcelas", "compra_fatura_id", "data_compra", "recorrente").
		Values(req.Nome, req.Descricao, req.LocalCompra, req.CategoriaID, req.ValorParcela, req.ParcelaAtual, req.QuantidadeParcelas, req.FaturaID, req.DataCompra, req.Recorrente).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).Scan(&req.ID); err != nil {
		return err
	}

	return
}
