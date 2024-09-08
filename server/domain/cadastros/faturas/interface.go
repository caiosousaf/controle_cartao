package faturas

import (
	"controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// IFatura define uma interface para os métodos de acesso a camada de dados
type IFatura interface {
	ListarFaturasCartao(p *utils.Parametros, id, usuarioID *uuid.UUID) (*faturas.FaturaPag, error)
	BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (*faturas.Fatura, error)
	BuscarFatura(idFatura, usuarioID *uuid.UUID) (*faturas.Fatura, error)
	ObterProximasFaturas(parcela_atual, qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error)
	VerificarFaturaCartao(data *string, idCartao, usuarioID *uuid.UUID) (*uuid.UUID, error)
	CadastrarFatura(req *faturas.Fatura) error
	AtualizarFatura(req *faturas.Fatura, idFatura *uuid.UUID) error
	AtualizarStatusFatura(req *faturas.Fatura, idFatura *uuid.UUID) error
	CartaoPertenceUsuario(idCartao, usuarioID *uuid.UUID) bool
}
