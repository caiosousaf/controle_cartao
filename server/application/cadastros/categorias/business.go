package categorias

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/categorias"
	infra "controle_cartao/infrastructure/cadastros/categorias"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// CadastrarCategoria contém a regra de negócio para cadastrar uma categoria
func CadastrarCategoria(req *ReqCategoria) (id *uuid.UUID, err error) {
	const (
		msgErrPadrao         = "Erro ao cadastrar nova categoria"
		msgErrPadraoListagem = "Erro ao listar categoria"
	)
	var (
		p utils.Parametros
	)

	var reqInfra = new(infra.Categorias)

	db, err := database.Conectar()
	if err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := categorias.NovoRepo(db)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}

	p.Filtros = make(map[string][]string)
	p.Filtros["nome_exato"] = []string{*req.Nome}
	p.Limite = 1
	lista, err := repo.ListarCategorias(&p)
	if err != nil {
		return id, utils.Wrap(err, msgErrPadraoListagem)
	}

	// Verifica se já existe algum cartão com esse nome
	if len(lista.Dados) > 0 {
		if lista.Dados[0].DataDesativacao != nil {

		} else {
			return id, utils.NewErr("Já existe uma categoria ativa com esse nome")
		}
	}

	if err = repo.CadastrarCategoria(reqInfra); err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}

	id = reqInfra.ID

	return
}

// AtualizarCategoria contém a regra de negócio para atualizar uma categoria
func AtualizarCategoria(req *ReqCategoria, idCategoria *uuid.UUID) (err error) {
	const (
		msgErrPadrao         = "Erro ao atualizar categoria"
		msgErrPadraoListagem = "Erro ao buscar categoria"
	)
	var (
		p utils.Parametros
	)

	var reqInfra = new(infra.Categorias)

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := categorias.NovoRepo(db)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	p.Filtros = make(map[string][]string)
	p.Filtros["nome_exato"] = []string{*req.Nome}
	p.Limite = 1
	lista, err := repo.ListarCategorias(&p)
	if err != nil {
		return utils.Wrap(err, msgErrPadraoListagem)
	}

	if len(lista.Dados) > 0 {
		if lista.Dados[0].DataDesativacao != nil {
			return utils.NewErr("Categoria desativada!")
		} else if *idCategoria == *lista.Dados[0].ID {
			return
		} else {
			return utils.NewErr("Já existe um cartão ativo com esse nome!")
		}
	}

	if err = repo.AtualizarCategoria(reqInfra, idCategoria); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}

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

// BuscarCategoria contém a regra de negócio para buscar uma categoria
func BuscarCategoria(idCategoria *uuid.UUID) (res ResCategorias, err error) {
	const msgErrPadrao = "Erro ao buscar categoria"

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := categorias.NovoRepo(db)

	resCategoria, err := repo.BuscarCategoria(idCategoria)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res = ResCategorias{
		ID:              resCategoria.ID,
		Nome:            resCategoria.Nome,
		DataCriacao:     resCategoria.DataCriacao,
		DataDesativacao: resCategoria.DataDesativacao,
	}

	return
}
