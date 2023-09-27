package faturas

import (
	"github.com/google/uuid"
	"time"
)

// Fatura estrutura para definição de modelo de fatura para uso na camada de dados
type Fatura struct {
	ID             *uuid.UUID `alias:"TFC" sql:"id" apelido:"id"`
	Nome           *string    `alias:"TFC" sql:"nome" apelido:"nome"`
	FaturaCartaoID *uuid.UUID `alias:"TFC" sql:"fatura_cartao_id" apelido:"cartao_id"`
	NomeCartao     *string    `alias:"TC" sql:"nome" apelido:"nome_cartao"`
	DataCriacao    *time.Time `alias:"TFC" sql:"data_criacao" apelido:"data_criacao"`
	DataVencimento *string    `alias:"TFC" sql:"data_vencimento" apelido:"data_vencimento"`
}

// FaturaPag estrutura para retorno de lista de dados paginada
type FaturaPag struct {
	Dados []Fatura
	Prox  *bool
	Total *int64
}
