package recorrente

import (
	model "controle_cartao/infrastructure/cadastros/compras/recorrente"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// IRecorrente define uma 'interface' para os métodos de acesso à camada de dados
type IRecorrente interface {
	ListarComprasRecorrentes(params *utils.Parametros, usuarioID *uuid.UUID) (*model.RecorrentesPag, error)
	ObterFaturaCartaoGeral(usuarioID *uuid.UUID) (*uuid.UUID, error)
	CadastrarCompraRecorrente(req *model.ComprasRecorrentes) (err error)
	ObterPrevisaoGastos(usuarioID *uuid.UUID) (gastos *model.PrevisaoGastosPag, err error)
}
