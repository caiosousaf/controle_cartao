package cartao

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/cartao"
	"controle_cartao/utils"
)

// ListarCartoes contém a regra de negócio para listar os cartões
func ListarCartoes(p *utils.Parametros) (res *ResPag, err error) {
	msgErrPadrao := "Erro ao listar cartões"

	res = new(ResPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := cartao.NovoRepo(db)

	listaCartoes, err := repo.ListarCartoes(p)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res.Dados = make([]Res, len(listaCartoes.Dados))
	for i := 0; i < len(listaCartoes.Dados); i++ {
		if err = utils.ConvertStructByAlias(&listaCartoes.Dados[i], &res.Dados[i]); err != nil {
			return res, utils.Wrap(err, msgErrPadrao)
		}
	}

	res.Total, res.Prox = listaCartoes.Total, listaCartoes.Prox

	return
}
