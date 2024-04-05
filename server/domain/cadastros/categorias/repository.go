package categorias

import (
	"controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/infrastructure/cadastros/categorias/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
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

// RemoverCategoria é um gerenciador de fluxo de dados para remover uma categoria no banco de dados
func (r *repo) RemoverCategoria(idCategoria *uuid.UUID) error {
	return r.Data.RemoverCategoria(idCategoria)
}

// ReativarCategoria é um gerenciador de fluxo de dados para reativar uma categoria no banco de dados
func (r *repo) ReativarCategoria(idCategoria *uuid.UUID) error {
	return r.Data.ReativarCategoria(idCategoria)
}
