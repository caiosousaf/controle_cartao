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

// CadastrarCategoria é um gerenciador de fluxo de dados para cadastrar uma nova categoria
func (r *repo) CadastrarCategoria(req *categorias.Categorias) error {
	return r.Data.CadastrarCategoria(req)
}

// AtualizarCategoria é um gerenciador de fluxo de dados para atualizar uma categoria
func (r *repo) AtualizarCategoria(req *categorias.Categorias, idCategoria *uuid.UUID) error {
	return r.Data.AtualizarCategoria(req, idCategoria)
}

// ListarCategorias é um gerenciador de fluxo de dados para listar as categorias
func (r *repo) ListarCategorias(params *utils.Parametros) (*categorias.CategoriasPag, error) {
	return r.Data.ListarCategorias(params)
}

// RemoverCategoria é um gerenciador de fluxo de dados para remover uma categoria no banco de dados
func (r *repo) RemoverCategoria(idCategoria, usuarioID *uuid.UUID) error {
	return r.Data.RemoverCategoria(idCategoria, usuarioID)
}

// ReativarCategoria é um gerenciador de fluxo de dados para reativar uma categoria no banco de dados
func (r *repo) ReativarCategoria(idCategoria, usuarioID *uuid.UUID) error {
	return r.Data.ReativarCategoria(idCategoria, usuarioID)
}

// BuscarCategoria é um gerenciador de fluxo de dados para buscar uma categoria no banco de dados a partir do seu id
func (r *repo) BuscarCategoria(idCategoria, usuarioID *uuid.UUID) (*categorias.Categorias, error) {
	return r.Data.BuscarCategoria(idCategoria, usuarioID)
}
