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
