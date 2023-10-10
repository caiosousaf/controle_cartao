package compras

import (
	"controle_cartao/infrastructure/cadastros/compras"
	"controle_cartao/utils"
)

// ICompra define uma interface para os métodos de acesso a camada de dados
type ICompra interface {
	CadastrarCompra(req *compras.Compras) error
	ListarCompras(params *utils.Parametros) (*compras.ComprasPag, error)
	ObterTotalComprasValor(params *utils.Parametros) (*compras.TotalComprasValor, error)
}
