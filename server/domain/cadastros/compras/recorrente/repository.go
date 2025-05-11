package recorrente

import (
	model "controle_cartao/infrastructure/cadastros/compras/recorrente"
	"controle_cartao/infrastructure/cadastros/compras/recorrente/postgres"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

type repo struct {
	Data *postgres.DBRecorrentes
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBRecorrentes{DB: novoDB},
	}
}

// ListarComprasRecorrentes é um gerenciador de fluxo de dados para listar as compras recorrentes
func (r *repo) ListarComprasRecorrentes(parmams *utils.Parametros, usuarioID *uuid.UUID) (res *model.RecorrentesPag, err error) {
	return r.Data.ListarComprasRecorrentes(parmams, usuarioID)
}
