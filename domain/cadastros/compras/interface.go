package compras

import "controle_cartao/infrastructure/cadastros/compras"

// ICompra define uma interface para os métodos de acesso a camada de dados
type ICompra interface {
	CadastrarCompra(req *compras.Compras) error
}
