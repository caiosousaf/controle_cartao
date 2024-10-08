package compras

import (
	"github.com/google/uuid"
	"time"
)

var (
	BaseURLCompras    = "/cadastros/compras"
	BaseURLComprasPdf = "/cadastros/compras/pdf"
)

// UserStateCompras Struct para armazenar o estado da conversa do usuário para ações de compras
type UserStateCompras struct {
	FaturaID          *string
	CurrentStep       *string
	CurrentStepBool   bool
	NovaCompraData    NovaCompra
	ObterTotalCompras ObterCompras
}

// NovaCompra Struct para armazenar os dados de uma nova compra
type NovaCompra struct {
	Nome               *string    `json:"nome,omitempty"`
	Descricao          *string    `json:"descricao,omitempty"`
	LocalCompra        *string    `json:"local_compra,omitempty"`
	CategoriaID        *uuid.UUID `json:"categoria_id,omitempty"`
	ValorParcela       *float64   `json:"valor_parcela,omitempty"`
	ParcelaAtual       *int64     `json:"parcela_atual,omitempty"`
	QuantidadeParcelas *int64     `json:"quantidade_parcelas,omitempty"`
	DataCompra         *string    `json:"data_compra,omitempty"`
}

// ObterCompras struct que armazena os dados de filtragem para obter valor compras
type ObterCompras struct {
	DataEspecifica   *string
	Pago             *bool
	UltimaParcela    *bool
	StepComprasTotal *bool
}

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
	DataCriacao        *time.Time `json:"data_criacao" apelido:"data_criacao"`
}

// ResComprasPag modela uma lista de respostas com suporte para paginação de compras
type ResComprasPag struct {
	Dados []ResCompras `json:"dados,omitempty"`
	Prox  *bool        `json:"prox,omitempty"`
	Total *int64       `json:"total,omitempty"`
}

// ResObterComprasTotal modela uma lista de resposta para obter o valor total das compras
type ResObterComprasTotal struct {
	Total *string `json:"total"`
}
