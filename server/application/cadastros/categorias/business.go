package categorias

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/categorias"
	"controle_cartao/utils"
)

// ListarCategorias contém a regra de negócio para listar as categorias
func ListarCategorias(params *utils.Parametros) (res *ResCategoriasPag, err error) {
	const msgErrPadrao = "Erro ao listar categorias"

	res = new(ResCategoriasPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}

	defer db.Close()

	repo := categorias.NovoRepo(db)

	listaCategorias, err := repo.ListarCategorias(params)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res.Dados = make([]ResCategorias, len(listaCategorias.Dados))
	for i := 0; i < len(listaCategorias.Dados); i++ {
		if err = utils.ConvertStructByAlias(&listaCategorias.Dados[i], &res.Dados[i]); err != nil {
			return res, utils.Wrap(err, msgErrPadrao)
		}
	}

	res.Total, res.Prox = listaCategorias.Total, listaCategorias.Prox

	return
}
