package cartao

import "database/sql"

func NovoRepo(DB *sql.DB) ICartao {
	return novoRepo(DB)
}
