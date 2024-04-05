package categorias

import (
	"controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/utils"
)

// ICategoria define uma interface para os métodos de acesso a camada de dados
type ICategoria interface {
	ListarCategorias(params *utils.Parametros) (*categorias.CategoriasPag, error)
}
