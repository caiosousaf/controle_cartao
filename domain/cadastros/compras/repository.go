package compras

import (
	"controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/infrastructure/cadastros/compras/postgres"
	"database/sql"
)

type repo struct {
	Data *postgres.DBCompra
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBCompra{DB: novoDB},
	}
}

// CadastrarCompra é um gerenciador de fluxo de dados para cadastrar uma nova compra no banco de dados
func (r *repo) CadastrarCompra(req *compras.Compras) error {
	return r.Data.CadastrarCompra(req)
}
