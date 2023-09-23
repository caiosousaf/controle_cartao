package cartao

import (
	"controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/infrastructure/cadastros/cartao/postgres"
	"controle_cartao/utils"
	"database/sql"
)

type repo struct {
	Data *postgres.DBCartao
}

func novoRepo(novoDB *sql.DB) *repo {
	return &repo{
		Data: &postgres.DBCartao{DB: novoDB},
	}
}

// ListarCartoes é um gerenciador de fluxo de dados para listar os cartões no banco de dados
func (r *repo) ListarCartoes(p *utils.Parametros) (*cartao.CartaoPag, error) {
	return r.Data.ListarCartoes(p)
}
