package recorrente

import (
	"github.com/google/uuid"
	"time"
)

// Recorrentes é a estrutura que define os dados das compras recorrentes
type Recorrentes struct {
	ID           *uuid.UUID `alias:"TCR" sql:"id" apelido:"id"`
	Nome         *string    `alias:"TCR" sql:"nome" apelido:"nome"`
	Descricao    *string    `alias:"TCR" sql:"descricao" apelido:"descricao"`
	CategoriaID  *uuid.UUID `alias:"TCR" sql:"compra_categoria_id" apelido:"categoria_id"`
	LocalCompra  *string    `alias:"TCR" sql:"local_compra" apelido:"local_compra"`
	ValorParcela *float64   `alias:"TCR" sql:"valor_parcela" apelido:"valor_parcela"`
	Ativo        *bool      `alias:"TCR" sql:"ativo" apelido:"ativo"`
	DataCriacao  *time.Time `alias:"TCR" sql:"data_criacao" apelido:"data_criacao"`
}

// RecorrentesPag é a estrutura de resposta paginada dos dados das compras recorrentes
type RecorrentesPag struct {
	Dados []Recorrentes
	Prox  *bool
	Total *int64
}

// ComprasRecorrentes é a estrutura que define os dados para cadastro das compras recorrentes
type ComprasRecorrentes struct {
	ID                 *uuid.UUID
	Nome               *string
	Descricao          *string
	LocalCompra        *string
	CategoriaID        *uuid.UUID
	ValorParcela       *float64
	ParcelaAtual       *int64
	QuantidadeParcelas *int64
	FaturaID           *uuid.UUID
	DataCompra         *string
	Recorrente         bool
}

// PrevisaoGastos é a estrutura que define os dados para previsão de gastos dos próximos meses
type PrevisaoGastos struct {
	MesAno string  `apelido:"mes_ano"`
	Valor  float64 `apelido:"valor"`
}

// PrevisaoGastosPag é a estrutura que define o modelo de retorno em forma de array dos gastos previstos
type PrevisaoGastosPag struct {
	Dados []PrevisaoGastos
}
