package categorias

import (
	"controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/infrastructure/cadastros/categorias/postgres"
	"controle_cartao/utils"
	"database/sql"
)

type repo struct {
	Data *postgres.DBCategoria
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBCategoria{DB: novoDB},
	}
}

// ListarCategorias é um gerenciador de fluxo de dados para listar as categorias
func (r *repo) ListarCategorias(params *utils.Parametros) (*categorias.CategoriasPag, error) {
	return r.Data.ListarCategorias(params)
}
