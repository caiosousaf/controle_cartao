package categorias

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/categorias"
	"controle_cartao/utils"
	"github.com/google/uuid"
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

// RemoverCategoria contém a regra de negócio para remover uma categoria
func RemoverCategoria(idCategoria *uuid.UUID) error {
	const msgErrPadrao = "Erro ao remover categoria"

	db, err := database.Conectar()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := categorias.NovoRepo(db)

	if err := repo.RemoverCategoria(idCategoria); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return err
}

// ReativarCategoria contém a regra de negócio para reativar uma categoria
func ReativarCategoria(idCategoria *uuid.UUID) error {
	const msgErrPadrao = "Erro ao reativar categoria"

	db, err := database.Conectar()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := categorias.NovoRepo(db)

	if err := repo.ReativarCategoria(idCategoria); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return err
}
