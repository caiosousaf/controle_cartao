package compras

import (
	"github.com/google/uuid"
	"time"
)

// Req modela uma requisição para a criação de uma compra
type Req struct {
	Nome               *string    `json:"nome" apelido:"nome"`
	Descricao          *string    `json:"descricao" apelido:"descricao"`
	LocalCompra        *string    `json:"local_compra" apelido:"local_compra"`
	CategoriaID        *uuid.UUID `json:"categoria_id" apelido:"categoria_id"`
	ValorParcela       *float64   `json:"valor_parcela" apelido:"valor_parcela"`
	ParcelaAtual       *int64     `json:"parcela_atual" apelido:"parcela_atual"`
	QuantidadeParcelas *int64     `json:"quantidade_parcelas" apelido:"quantidade_parcelas"`
	FaturaID           *uuid.UUID `json:"fatura_id" apelido:"fatura_id"`
	AgrupamentoID      *uuid.UUID `json:"-" apelido:"agrupamento_id"`
	DataCompra         *string    `json:"data_compra" apelido:"data_compra"`
}

var (
	// colunasFaturasPdf é a variavel que define as colunas que serão usadas no retorno do pdf da tabela de faturas
	colunasFaturasPdf = []string{"Nome", "Local Compra", "Categoria", "Valor Parcela", "Parcela Atual", "Quantidade Parcelas", "Data Compra"}
	// colunasMesesFaturasCartao é a variavel que define as colunas que serão usadas no retorno do pdf para a tabela com os meses das faturas de um cartão
	colunasMesesFaturasCartao = []string{"Nome", "Status", "Total"}
)

// ResCompras modela uma resposta para listagem e busca de compras
type ResCompras struct {
	ID                 *uuid.UUID `json:"id" apelido:"id"`
	Nome               *string    `json:"nome" apelido:"nome"`
	Descricao          *string    `json:"descricao" apelido:"descricao"`
	LocalCompra        *string    `json:"local_compra" apelido:"local_compra"`
	CategoriaID        *uuid.UUID `json:"categoria_id" apelido:"categoria_id"`
	CategoriaNome      *string    `json:"categoria_nome" apelido:"categoria_nome"`
	ValorParcela       *float64   `json:"valor_parcela" apelido:"valor_parcela"`
	ParcelaAtual       *int64     `json:"parcela_atual" apelido:"parcela_atual"`
	QuantidadeParcelas *int64     `json:"quantidade_parcelas" apelido:"quantidade_parcelas"`
	FaturaID           *uuid.UUID `json:"fatura_id" apelido:"fatura_id"`
	NomeFatura         *string    `json:"nome_fatura" apelido:"fatura_nome"`
	DataCompra         *string    `json:"data_compra" apelido:"data_compra"`
	Recorrente         *bool      `json:"recorrente" apelido:"recorrente"`
	AgrupamentoID      *uuid.UUID `json:"agrupamento_id,omitempty" apelido:"agrupamento_id"`
	DataCriacao        *time.Time `json:"data_criacao" apelido:"data_criacao"`
}

// ResComprasPag modela uma lista de respostas com suporte para paginação de compras
type ResComprasPag struct {
	Dados []ResCompras `json:"dados,omitempty"`
	Prox  *bool        `json:"prox,omitempty"`
	Total *int64       `json:"total,omitempty"`
}

// ResTotalComprasValor modela uma estrutura para obter o valor total das compras
type ResTotalComprasValor struct {
	Total *string `json:"total" apelido:"valor_total"`
}
