package compras

import (
	"controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/utils"

	"github.com/google/uuid"
)

// ICompra define uma ‘interface’ para os métodos de acesso à camada de dados
type ICompra interface {
	CadastrarCompra(req *compras.Compras) error
	ListarCompras(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.ComprasPag, error)
	ObterTotalComprasValor(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.TotalComprasValor, error)
	AtualizarCompra(req *compras.Compras, usuarioID, compraID *uuid.UUID, recorrente, atualizarTodasParcelas bool) error
	RemoverCompra(compraID, usuarioID *uuid.UUID, recorrente, removerTodasParcelas bool) error
	VerificaCompraRecorrente(compraID *uuid.UUID) (recorrente *bool, err error)
}
