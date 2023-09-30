package faturas

import "database/sql"

func NovoRepo(DB *sql.DB) IFatura {
	return novoRepo(DB)
}
