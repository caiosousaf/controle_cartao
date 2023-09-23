package postgres

import (
	model "controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/utils"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// DBCartao é uma estrutura base para acesso aos metodos do banco postgres para manipulação de cartões
type DBCartao struct {
	DB *sql.DB
}

// ListarCartoes lista de forma paginada os dados dos cartões cadastrados no banco de dados
func (pg *DBCartao) ListarCartoes(p *utils.Parametros) (res *model.CartaoPag, err error) {
	var t model.Cartao

	res = new(model.CartaoPag)

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select("id, nome, data_criacao, data_desativacao").
		From("public.t_cartao")

	dados, prox, total, err := utils.ConfigurarPaginacao(p, &t, &consultaSql)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Cartao), prox, total

	return
}

// BuscarCartao busca os dados de um cartão no banco de dados dado o seu id
func (pg *DBCartao) BuscarCartao(id *uuid.UUID) (res *model.Cartao, err error) {
	res = new(model.Cartao)
	if err = sq.StatementBuilder.RunWith(pg.DB).Select("id, nome, data_criacao, data_desativacao").
		From("public.t_cartao").
		Where("id = $1", id).
		Where("data_desativacao ISNULL").Scan(&res.ID, &res.Nome, &res.DataCriacao, &res.DataDesativacao); err != nil {
		return res, err
	}

	return
}
