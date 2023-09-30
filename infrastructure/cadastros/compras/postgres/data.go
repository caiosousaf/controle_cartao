package postgres

import (
	model "controle_cartao/infrastructure/cadastros/compras"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

// DBCompra é uma estrutura para acesso aos metodos do banco postgres para manipulação das compras
type DBCompra struct {
	DB *sql.DB
}

// CadastrarCompras cadastra uma nova compra no banco de dados
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
