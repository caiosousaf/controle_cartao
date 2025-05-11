package recorrente

import (
	"github.com/google/uuid"
	"time"
)

// ResRecorrentes modela uma resposta para listagem de compras recorrentes
type ResRecorrentes struct {
	ID           *uuid.UUID `json:"id" apelido:"id"`
	Nome         *string    `json:"nome" apelido:"nome"`
	Descricao    *string    `json:"descricao" apelido:"descricao"`
	CategoriaID  *uuid.UUID `json:"compra_categoria_id" apelido:"categoria_id"`
	LocalCompra  *string    `json:"local_compra" apelido:"local_compra"`
	ValorParcela *float64   `json:"valor_parcela" apelido:"valor_parcela"`
	Ativo        *bool      `json:"ativo" apelido:"ativo"`
	DataCriacao  *time.Time `json:"data_criacao" apelido:"data_criacao"`
}

// ResRecorrentesPag modela uma lista de respostas paginada de compras recorrentes
type ResRecorrentesPag struct {
	Dados []ResRecorrentes `json:"dados,omitempty"`
	Prox  *bool            `json:"prox,omitempty"`
	Total *int64           `json:"total,omitempty"`
}
