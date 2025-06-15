package compras

import (
	"github.com/google/uuid"
	"time"
)

// TotalComprasValor estrutura para definição de total de compras para uso na camada de dados
type TotalComprasValor struct {
	Total *string `alias:"TCF" sql:"valor_parcela" apelido:"valor_total"`
}

// Compras estrutura para definição de modelo de compra para uso na camada de dados
type Compras struct {
	ID                 *uuid.UUID `alias:"TCF" sql:"id" apelido:"id"`
	Nome               *string    `alias:"TCF" sql:"nome" apelido:"nome"`
	Descricao          *string    `alias:"TCF" sql:"descricao" apelido:"descricao"`
	LocalCompra        *string    `alias:"TCF" sql:"local_compra" apelido:"local_compra"`
	CategoriaID        *uuid.UUID `alias:"TCF" sql:"compra_categoria_id" apelido:"categoria_id"`
	CategoriaNome      *string    `alias:"TCC" sql:"nome" apelido:"categoria_nome"`
	ValorParcela       *float64   `alias:"TCF" sql:"valor_parcela" apelido:"valor_parcela"`
	ParcelaAtual       *int64     `alias:"TCF" sql:"parcela_atual" apelido:"parcela_atual"`
	QuantidadeParcelas *int64     `alias:"TCF" sql:"qtd_parcelas" apelido:"quantidade_parcelas"`
	FaturaID           *uuid.UUID `alias:"TCF" sql:"compra_fatura_id" apelido:"fatura_id"`
	NomeFatura         *string    `alias:"TFC" sql:"nome" apelido:"fatura_nome"`
	DataCompra         *string    `alias:"TCF" sql:"data_compra" apelido:"data_compra"`
	Recorrente         *bool      `alias:"TCF" sql:"recorrente"`
	AgrupamentoID      *uuid.UUID `alias:"TCF" sql:"agrupamento_id" apelido:"agrupamento_id"`
	DataCriacao        *time.Time `alias:"TCF" sql:"data_criacao" apelido:"data_criacao"`
}

// ComprasPag estrutura para retorno de lista de dados paginada
type ComprasPag struct {
	Dados []Compras
	Prox  *bool
	Total *int64
}
