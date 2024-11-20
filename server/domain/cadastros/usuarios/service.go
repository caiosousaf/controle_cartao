package usuarios

import "database/sql"

func NovoRepo(DB *sql.DB) IUsuario {
	return novoRepo(DB)
}
