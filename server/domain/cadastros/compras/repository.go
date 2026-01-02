package compras

import (
	"controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/infrastructure/cadastros/compras/postgres"
	"controle_cartao/utils"
	"database/sql"

	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBCompra
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBCompra{DB: novoDB},
	}
}

// CadastrarCompra é um gerenciador de fluxo de dados para cadastrar uma nova compra no banco de dados
func (r *repo) CadastrarCompra(req *compras.Compras) error {
	return r.Data.CadastrarCompra(req)
}

// ListarCompras é um gerenciador de fluxo de dados para listar as compras no banco de dados
func (r *repo) ListarCompras(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.ComprasPag, error) {
	return r.Data.ListarCompras(params, usuarioID)
}

// ObterTotalComprasValor é um gerenciador de fluxo de dados para obter o valor total de compras dado os devidos filtros
func (r *repo) ObterTotalComprasValor(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.TotalComprasValor, error) {
	return r.Data.ObterTotalComprasValor(params, usuarioID)
}

// AtualizarCompra é um gerenciador de fluxo de dados para atualizar compras
func (r *repo) AtualizarCompra(req *compras.Compras, usuarioID, compraID *uuid.UUID, recorrente, atualizarTodasParcelas bool) error {
	return r.Data.AtualizarCompra(req, usuarioID, compraID, recorrente, atualizarTodasParcelas)
}

// RemoverCompra é um gerenciador de fluxo de dados para remover uma compra
func (r *repo) RemoverCompra(compraID, usuarioID *uuid.UUID, recorrente, removerTodasParcelas bool) error {
	return r.Data.RemoverCompra(compraID, usuarioID, recorrente, removerTodasParcelas)
}

// VerificaCompraRecorrente é um gerenciador de fluxo de dados para verificar compra recorrente ou não
func (r *repo) VerificaCompraRecorrente(compraID *uuid.UUID) (recorrente *bool, err error) {
	return r.Data.VerificaCompraRecorrente(compraID)
}

// AnteciparParcelas é um gerenciador de fluxo de dados para antecipação de parcelas de uma compra especifica para uma fatura
func (r *repo) AnteciparParcelas(req *compras.ReqAntecipacaoParcelas, faturaID, usuarioID *uuid.UUID) error {
	return r.Data.AnteciparParcelas(req, faturaID, usuarioID)
}

// ObterParcelasDisponiveisAntecipacao é um gerenciador de fluxo de dados para obter as parcelas disponiveis para antecipação
func (r *repo) ObterParcelasDisponiveisAntecipacao(identificadorCompra, faturaID, usuarioID *uuid.UUID) ([]int64, error) {
	return r.Data.ObterParcelasDisponiveisAntecipacao(identificadorCompra, faturaID, usuarioID)
}
