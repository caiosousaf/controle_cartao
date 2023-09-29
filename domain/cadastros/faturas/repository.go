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
func (r *repo) ListarFaturasCartao(p *utils.Parametros, id *uuid.UUID) (*faturas.FaturaPag, error) {
	return r.Data.ListarFaturasCartao(p, id)
}

// BuscarFaturaCartao é um gerenciador de fluxo de dados para buscar a fatura de um cartão no banco de dados
func (r *repo) BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (*faturas.Fatura, error) {
	return r.Data.BuscarFaturaCartao(idCartao, idFatura)
}

// ObterProximasFaturas é um gerenciador de fluxo de dados para obter as próximas faturas de um cartão no banco de dados
func (r *repo) ObterProximasFaturas(qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error) {
	return r.Data.ObterProximasFaturas(qtd_parcelas, idFatura)
}

// VerificarFaturaCartao é um gerenciador de fluxo de dados para verificar se existe uma fatura de cartão para a data selecionada no banco de dados
func (r *repo) VerificarFaturaCartao(data *string, idCartao *uuid.UUID) (*uuid.UUID, error) {
	return r.Data.VerificarFaturaCartao(data, idCartao)
}

// CadastrarFatura é um gerenciador de fluxo de dados para cadastrar uma nova fatura no banco de dados
func (r *repo) CadastrarFatura(req *faturas.Fatura) error {
	return r.Data.CadastrarFatura(req)
}

// AtualizarFatura é um gerenciador de fluxo de dados para atualizar uma fatura no banco de dados
func (r *repo) AtualizarFatura(req *faturas.Fatura, idFatura *uuid.UUID) error {
	return r.Data.AtualizarFatura(req, idFatura)
}
