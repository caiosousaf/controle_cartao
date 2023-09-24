package cartao

import (
	"controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/infrastructure/cadastros/cartao/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBCartao
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBCartao{DB: novoDB},
	}
}

// CadastrarCartao é um gerenciador de fluxo de dados para cadastrar um novo cartão no banco de dados
func (r *repo) CadastrarCartao(req *cartao.Cartao) error {
	return r.Data.CadastrarCartao(req)
}

// ListarCartoes é um gerenciador de fluxo de dados para listar os cartões no banco de dados
func (r *repo) ListarCartoes(p *utils.Parametros) (*cartao.CartaoPag, error) {
	return r.Data.ListarCartoes(p)
}

// BuscarCartao é um gerenciador de fluxo de dados para buscar um cartão no banco de dados
func (r *repo) BuscarCartao(id *uuid.UUID) (*cartao.Cartao, error) {
	return r.Data.BuscarCartao(id)
}

// AtualizarCartao é um gerenciador de fluxo de dados para atualizar um cartão no banco de dados
func (r *repo) AtualizarCartao(req *cartao.Cartao, id *uuid.UUID) error {
	return r.Data.AtualizarCartao(req, id)
}

// RemoverCartao é um gerenciador de fluxo de dados para desativar um cartão no banco de dados
func (r *repo) RemoverCartao(id *uuid.UUID) error {
	return r.Data.RemoverCartao(id)
}
