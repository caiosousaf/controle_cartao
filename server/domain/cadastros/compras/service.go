package compras

import "database/sql"

func NovoRepo(DB *sql.DB) ICompra {
	return novoRepo(DB)
}
