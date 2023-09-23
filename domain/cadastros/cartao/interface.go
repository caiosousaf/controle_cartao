package cartao

import (
	"controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// ICartao define uma interface para os metodos de acesso a camada de dados
type ICartao interface {
	ListarCartoes(p *utils.Parametros) (*cartao.CartaoPag, error)
	BuscarCartao(id *uuid.UUID) (*cartao.Cartao, error)
}
