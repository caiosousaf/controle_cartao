package compras

import (
	"controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/utils"

	"github.com/google/uuid"
)

// ICompra define uma interface para os métodos de acesso a camada de dados
type ICompra interface {
	CadastrarCompra(req *compras.Compras) error
	ListarCompras(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.ComprasPag, error)
	ObterTotalComprasValor(params *utils.Parametros, usuarioID *uuid.UUID) (*compras.TotalComprasValor, error)
}
