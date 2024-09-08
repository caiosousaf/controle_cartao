package faturas

import (
	"controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/infrastructure/cadastros/faturas/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBFatura
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBFatura{DB: novoDB},
	}
}

// ListarFaturasCartao é um gerenciador de fluxo de dados para listar as faturas de um cartão no banco de dados
func (r *repo) ListarFaturasCartao(p *utils.Parametros, id, usuarioID *uuid.UUID) (*faturas.FaturaPag, error) {
	return r.Data.ListarFaturasCartao(p, id, usuarioID)
}

// BuscarFaturaCartao é um gerenciador de fluxo de dados para buscar a fatura de um cartão no banco de dados
func (r *repo) BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (*faturas.Fatura, error) {
	return r.Data.BuscarFaturaCartao(idCartao, idFatura)
}

// BuscarFatura é um gerenciador de fluxo de dados para buscar a fatura de um cartão no banco de dados
func (r *repo) BuscarFatura(idFatura, usuarioID *uuid.UUID) (*faturas.Fatura, error) {
	return r.Data.BuscarFatura(idFatura, usuarioID)
}

// ObterProximasFaturas é um gerenciador de fluxo de dados para obter as próximas faturas de um cartão no banco de dados
func (r *repo) ObterProximasFaturas(parcela_atual, qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error) {
	return r.Data.ObterProximasFaturas(parcela_atual, qtd_parcelas, idFatura)
}

// VerificarFaturaCartao é um gerenciador de fluxo de dados para verificar se existe uma fatura de cartão para a data selecionada no banco de dados
func (r *repo) VerificarFaturaCartao(data *string, idCartao, usuarioID *uuid.UUID) (*uuid.UUID, error) {
	return r.Data.VerificarFaturaCartao(data, idCartao, usuarioID)
}

// CadastrarFatura é um gerenciador de fluxo de dados para cadastrar uma nova fatura no banco de dados
func (r *repo) CadastrarFatura(req *faturas.Fatura) error {
	return r.Data.CadastrarFatura(req)
}

// AtualizarFatura é um gerenciador de fluxo de dados para atualizar uma fatura no banco de dados
func (r *repo) AtualizarFatura(req *faturas.Fatura, idFatura *uuid.UUID) error {
	return r.Data.AtualizarFatura(req, idFatura)
}

// AtualizarStatusFatura é um gerenciador de fluxo de dados para atualizar o status de uma fatura no banco de dados
func (r *repo) AtualizarStatusFatura(req *faturas.Fatura, idFatura *uuid.UUID) error {
	return r.Data.AtualizarStatusFatura(req, idFatura)
}

// CartaoPertenceUsuario é um gerenciador de fluxo de dados para validar se o cartão pertence ao usuário
func (r *repo) CartaoPertenceUsuario(idCartao, usuarioID *uuid.UUID) bool {
	return r.Data.CartaoPertenceUsuario(idCartao, usuarioID)
}
