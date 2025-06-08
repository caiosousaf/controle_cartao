package recorrente

import "database/sql"

func NovoRepo(DB *sql.DB) IRecorrente {
	return novoRepo(DB)
}
