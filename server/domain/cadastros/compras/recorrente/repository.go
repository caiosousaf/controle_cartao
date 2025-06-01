package recorrente

import (
	model "controle_cartao/infrastructure/cadastros/compras/recorrente"
	"controle_cartao/infrastructure/cadastros/compras/recorrente/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBRecorrentes
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBRecorrentes{DB: novoDB},
	}
}

// ListarComprasRecorrentes é um gerenciador de fluxo de dados para listar as compras recorrentes
func (r *repo) ListarComprasRecorrentes(parmams *utils.Parametros, usuarioID *uuid.UUID) (res *model.RecorrentesPag, err error) {
	return r.Data.ListarComprasRecorrentes(parmams, usuarioID)
}

// ObterFaturaCartaoGeral é um gerenciador de fluxo de dados para obter o id da fatura do cartão geral
func (r *repo) ObterFaturaCartaoGeral(usuarioID *uuid.UUID) (faturaID *uuid.UUID, err error) {
	return r.Data.ObterFaturaCartaoGeral(usuarioID)
}

// CadastrarCompraRecorrente é um gerenciador de fluxo para cadastrar uma compra recorrente
func (r *repo) CadastrarCompraRecorrente(req *model.ComprasRecorrentes) (err error) {
	return r.Data.CadastrarCompraRecorrente(req)
}

// ObterPrevisaoGastos é um gerenciador de fluxo para obter uma previsão de gastos
func (r *repo) ObterPrevisaoGastos(usuarioID *uuid.UUID) (gastos *model.PrevisaoGastosPag, err error) {
	return r.Data.ObterPrevisaoGastos(usuarioID)
}

// CadastrarNovaCompraRecorrente é um gerenciador de fluxo para cadastrar uma nova compra recorrente
func (r *repo) CadastrarNovaCompraRecorrente(req *model.Recorrentes, usuarioID *uuid.UUID) (err error) {
	return r.Data.CadastrarNovaCompraRecorrente(req, usuarioID)
}

// AtualizarCompraRecorrente é um gerenciador de fluxo para atualizar compra recorrente
func (r *repo) AtualizarCompraRecorrente(req *model.Recorrentes, usuarioID *uuid.UUID) (err error) {
	return r.Data.AtualizarCompraRecorrente(req, usuarioID)
}

// DesativarCompraRecorrente é um gerenciador de fluxo para desativar uma compra recorrente
func (r *repo) DesativarCompraRecorrente(recorrenteID, usuarioID *uuid.UUID) error {
	return r.Data.DesativarCompraRecorrente(recorrenteID, usuarioID)
}

// ReativarCompraRecorrente é um gerenciador de fluxo para reativar uma compra recorrente
func (r *repo) ReativarCompraRecorrente(recorrenteID, usuarioID *uuid.UUID) error {
	return r.Data.ReativarCompraRecorrente(recorrenteID, usuarioID)
}

// RemoverCompraRecorrente é um gerenciador de fluxo para remover uma compra recorrente
func (r *repo) RemoverCompraRecorrente(recorrenteID, usuarioID *uuid.UUID) error {
	return r.Data.RemoverCompraRecorrente(recorrenteID, usuarioID)
}
