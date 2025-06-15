package postgres

import (
	model "controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBCompra é uma estrutura para acesso aos metodos do banco postgres para manipulação das compras
type DBCompra struct {
	DB *sql.DB
}

// CadastrarCompra cadastra uma nova compra no banco de dados
func (pg *DBCompra) CadastrarCompra(req *model.Compras) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_compras_fatura").
		Columns("nome", "descricao", "local_compra", "compra_categoria_id", "valor_parcela", "parcela_atual", "qtd_parcelas", "agrupamento_id", "compra_fatura_id", "data_compra").
		Values(req.Nome, req.Descricao, req.LocalCompra, req.CategoriaID, req.ValorParcela, req.ParcelaAtual, req.QuantidadeParcelas, req.AgrupamentoID, req.FaturaID, req.DataCompra).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).Scan(&req.ID); err != nil {
		return err
	}

	return
}

// ListarCompras lista de forma paginada os dados das compras no banco de dados
func (pg *DBCompra) ListarCompras(params *utils.Parametros, usuarioID *uuid.UUID) (res *model.ComprasPag, err error) {
	var t model.Compras

	res = new(model.ComprasPag)

	campos, _, err := params.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_compras_fatura TCF").
		Join("public.t_categoria_compra TCC on TCC.id = TCF.compra_categoria_id").
		Join("public.t_fatura_cartao TFC on TFC.id = TCF.compra_fatura_id").
		Join("public.t_cartao TC ON TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TC.usuario_id": usuarioID,
		})

	consultaComFiltro := params.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"fatura_id":       utils.CriarFiltros("TFC.id = ?::UUID", utils.FlagFiltroEq),
		"categoria_id":    utils.CriarFiltros("TCC.id = ?::UUID", utils.FlagFiltroEq),
		"data_especifica": utils.CriarFiltros("TO_CHAR(TFC.data_vencimento, 'MM/YYYY') = ?", utils.FlagFiltroEq),
		"cartao_id":       utils.CriarFiltros("TC.id = ?::UUID", utils.FlagFiltroEq),
	}).
		PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(params, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Compras), prox, total

	return
}

// ObterTotalComprasValor obtém o total de compras dado os seus filtros no banco de dados
func (pg *DBCompra) ObterTotalComprasValor(params *utils.Parametros, usuarioID *uuid.UUID) (res *model.TotalComprasValor, err error) {
	res = new(model.TotalComprasValor)

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).
		Select("ROUND(COALESCE(CAST(SUM(valor_parcela) AS numeric), 0), 2) AS valor_total").
		From("public.t_compras_fatura TCF").
		Join("public.t_fatura_cartao TFC ON TFC.id = TCF.compra_fatura_id").
		Join("public.t_cartao TC ON TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TC.usuario_id": usuarioID,
		})

	consultaComFiltro := params.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"cartao_id":       utils.CriarFiltros("TC.id = ?::UUID", utils.FlagFiltroEq),
		"fatura_id":       utils.CriarFiltros("TFC.id = ?::UUID", utils.FlagFiltroEq),
		"data_especifica": utils.CriarFiltros("TO_CHAR(TFC.data_vencimento, 'MM/YYYY') = ?", utils.FlagFiltroEq),
		"ultima_parcela":  utils.CriarFiltros("(TCF.parcela_atual = TCF.qtd_parcelas) = ?::BOOLEAN", utils.FlagFiltroEq),
		"pago":            utils.CriarFiltros("(TFC.status <> 'Pago') = ?::BOOLEAN", utils.FlagFiltroEq),
	}).
		PlaceholderFormat(sq.Dollar)

	if err = consultaComFiltro.
		Scan(&res.Total); err != nil {
		return res, err
	}

	return
}

// AtualizarCompra atualiza todas as compras de um agrupamento
func (pg *DBCompra) AtualizarCompra(req *model.Compras, usuarioID, compraID *uuid.UUID, recorrente, atualizarTodasParcelas bool) error {
	queryAgrupamento := `
		agrupamento_id IN (
			SELECT TCF.agrupamento_id
			FROM t_compras_fatura TCF
			JOIN t_fatura_cartao TFC ON TFC.id = TCF.compra_fatura_id
			JOIN t_cartao TC ON TC.id = TFC.fatura_cartao_id
			WHERE TC.usuario_id = ?
			AND TCF.id = ?)
	`

	query := sq.
		Update("public.t_compras_fatura").
		Set("nome", req.Nome).
		Set("descricao", req.Descricao).
		Set("local_compra", req.LocalCompra).
		Set("compra_categoria_id", req.CategoriaID).
		Set("valor_parcela", req.ValorParcela).
		Set("data_compra", req.DataCompra).
		PlaceholderFormat(sq.Dollar).
		RunWith(pg.DB)

	if recorrente || !atualizarTodasParcelas {
		query = query.Where(sq.Eq{"id": compraID})
	} else {
		query = query.Where(sq.Expr(queryAgrupamento, usuarioID.String(), compraID.String()))
	}

	row, err := query.Exec()
	if err != nil {
		return err
	}

	affected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return utils.NewErr("compra não encontrada")
	}

	return err
}

// RemoverCompra remove todas as compras de um agrupamento
func (pg *DBCompra) RemoverCompra(compraID, usuarioID *uuid.UUID, recorrente, removerTodasParcelas bool) error {
	queryAgrupamento := `agrupamento_id IN (
			SELECT TCF.agrupamento_id
			FROM t_compras_fatura TCF
			JOIN t_fatura_cartao TFC ON TFC.id = TCF.compra_fatura_id
			JOIN t_cartao TC ON TC.id = TFC.fatura_cartao_id
			WHERE TC.usuario_id = ?
			AND TCF.id = ?)`

	query := sq.
		Delete("public.t_compras_fatura").
		PlaceholderFormat(sq.Dollar).
		RunWith(pg.DB)

	if recorrente || !removerTodasParcelas {
		query = query.Where(sq.Eq{"id": compraID})
	} else {
		query = query.Where(sq.Expr(queryAgrupamento, usuarioID.String(), compraID.String()))
	}

	row, err := query.Exec()
	if err != nil {
		return err
	}

	affected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return utils.NewErr("compra não encontrada")
	}

	return err
}

// VerificaCompraRecorrente verifica se a compra é recorrente ou não
func (pg *DBCompra) VerificaCompraRecorrente(compraID *uuid.UUID) (recorrente *bool, err error) {
	query := sq.Select("recorrente").
		From("t_compras_fatura").
		Where(sq.Eq{"id": compraID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(pg.DB).
		QueryRow()

	if err = query.Scan(&recorrente); err != nil {
		return recorrente, err
	}

	return
}
