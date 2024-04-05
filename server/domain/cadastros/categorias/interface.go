package categorias

import (
	"controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// ICategoria define uma interface para os métodos de acesso a camada de dados
type ICategoria interface {
	CadastrarCategoria(req *categorias.Categorias) error
	AtualizarCategoria(req *categorias.Categorias, idCategoria *uuid.UUID) error
	ListarCategorias(params *utils.Parametros) (*categorias.CategoriasPag, error)
	RemoverCategoria(idCategoria *uuid.UUID) error
	ReativarCategoria(idCategoria *uuid.UUID) error
	BuscarCategoria(idCategoria *uuid.UUID) (*categorias.Categorias, error)
}
