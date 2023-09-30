package compras

import (
	"github.com/google/uuid"
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
	DataCompra         *string    `json:"data_compra" apelido:"data_compra"`
}
