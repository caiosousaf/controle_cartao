package faturas

import (
	"controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// IFatura define uma interface para os métodos de acesso a camada de dados
type IFatura interface {
	ListarFaturasCartao(p *utils.Parametros, id *uuid.UUID) (*faturas.FaturaPag, error)
	BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (*faturas.Fatura, error)
	BuscarFatura(idFatura *uuid.UUID) (*faturas.Fatura, error)
	ObterProximasFaturas(qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error)
	VerificarFaturaCartao(data *string, idCartao *uuid.UUID) (*uuid.UUID, error)
	CadastrarFatura(req *faturas.Fatura) error
	AtualizarFatura(req *faturas.Fatura, idFatura *uuid.UUID) error
	AtualizarStatusFatura(req *faturas.Fatura, idFatura *uuid.UUID) error
}
