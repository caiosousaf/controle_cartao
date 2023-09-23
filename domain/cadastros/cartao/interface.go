package cartao

import (
	"controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/utils"
)

// ICartao define uma interface para os metodos de acesso a camada de dados
type ICartao interface {
	ListarCartoes(p *utils.Parametros) (*cartao.CartaoPag, error)
}
