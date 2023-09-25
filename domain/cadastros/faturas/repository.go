package faturas

import (
	"controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/infrastructure/cadastros/faturas/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBFatura
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBFatura{DB: novoDB},
	}
}

// ListarFaturasCartao é um gerenciador de fluxo de dados para listar as faturas de um cartão no banco de dados
func (r *repo) ListarFaturasCartao(p *utils.Parametros, id *uuid.UUID) (*faturas.FaturaPag, error) {
	return r.Data.ListarFaturasCartao(p, id)
}
