package postgres

import (
	model "controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

// DBCompra é uma estrutura para acesso aos metodos do banco postgres para manipulação das compras
type DBCompra struct {
	DB *sql.DB
}

// CadastrarCompra cadastra uma nova compra no banco de dados
func (pg *DBCompra) CadastrarCompra(req *model.Compras) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_compras_fatura").
		Columns("nome", "descricao", "local_compra", "compra_categoria_id", "valor_parcela", "parcela_atual", "qtd_parcelas", "compra_fatura_id", "data_compra").
		Values(req.Nome, req.Descricao, req.LocalCompra, req.CategoriaID, req.ValorParcela, req.ParcelaAtual, req.QuantidadeParcelas, req.FaturaID, req.DataCompra).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).Scan(&req.ID); err != nil {
		return err
	}

	return
}

// ListarCompras lista de forma paginada os dados das compras no banco de dados
func (pg *DBCompra) ListarCompras(params *utils.Parametros) (res *model.ComprasPag, err error) {
	var t model.Compras

	res = new(model.ComprasPag)

	campos, _, err := params.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_compras_fatura TCF").
		Join("public.t_categoria_compra TCC on TCC.id = TCF.compra_categoria_id").
		Join("public.t_fatura_cartao TFC on TFC.id = TCF.compra_fatura_id")

	consultaComFiltro := params.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"fatura_id": utils.CriarFiltros("TFC.id = ?::UUID", utils.FlagFiltroEq),
	}).
		PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(params, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Compras), prox, total

	return
}
