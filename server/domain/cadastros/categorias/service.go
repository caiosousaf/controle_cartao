package categorias

import "database/sql"

func NovoRepo(DB *sql.DB) ICategoria {
	return novoRepo(DB)
}
