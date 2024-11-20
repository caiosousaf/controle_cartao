package cartao

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/cartao"
	infra "controle_cartao/infrastructure/cadastros/cartao"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// CadastrarCartao contém a regra de negócio para cadastrar um novo cartão
func CadastrarCartao(req *Req) (id *uuid.UUID, err error) {
	const (
		msgErrPadrao         = "Erro ao cadastrar novo cartão"
		msgErrPadraoListagem = "Erro ao listar cartão por nome"
	)
	var (
		p utils.Parametros
	)

	var reqInfra = new(infra.Cartao)

	db, err := database.Conectar()
	if err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := cartao.NovoRepo(db)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}

	p.Filtros = make(map[string][]string)
	p.Filtros["nome_exato"] = []string{*req.Nome}
	p.Filtros["usuario_id"] = []string{req.UsuarioID.String()}
	p.Limite = 1
	lista, err := repo.ListarCartoes(&p)
	if err != nil {
		return id, utils.Wrap(err, msgErrPadraoListagem)
	}

	// Verifica se já existe algum cartão com esse nome
	if len(lista.Dados) > 0 {
		if lista.Dados[0].DataDesativacao != nil {

		} else {
			return id, utils.NewErr("Já existe um cartão ativo com esse nome")
		}
	}

	if err = repo.CadastrarCartao(reqInfra); err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}

	id = reqInfra.ID

	return
}

// ListarCartoes contém a regra de negócio para listar os cartões
func ListarCartoes(p *utils.Parametros, usuarioID *uuid.UUID) (res *ResPag, err error) {
	msgErrPadrao := "Erro ao listar cartões"

	res = new(ResPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := cartao.NovoRepo(db)

	p.AdicionarFiltro("usuario_id", usuarioID.String())

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

// BuscarCartao contém a regra de negócio para buscar um cartão
func BuscarCartao(id, usuarioID *uuid.UUID) (res *Res, err error) {
	const msgErrPadrao = "Erro ao buscar um cartão"

	res = new(Res)

	db, err := database.Conectar()
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	defer db.Close()

	repo := cartao.NovoRepo(db)

	buscaCartao, err := repo.BuscarCartao(id, usuarioID)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res = &Res{
		ID:              buscaCartao.ID,
		Nome:            buscaCartao.Nome,
		UsuarioID:       buscaCartao.UsuarioID,
		DataCriacao:     buscaCartao.DataCriacao,
		DataDesativacao: buscaCartao.DataDesativacao,
	}

	return
}

// AtualizarCartao contém a regra de negócio para atualizar um cartão
func AtualizarCartao(req *ReqAtualizar, id, usuarioID *uuid.UUID) (err error) {
	const (
		msgErrPadrao         = "Erro ao atualizar cartão"
		msgErrPadraoListagem = "Erro ao listar cartão por nome"
	)
	var (
		p utils.Parametros
	)

	var reqInfra = new(infra.Cartao)

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := cartao.NovoRepo(db)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	p.Filtros = make(map[string][]string)
	p.Filtros["nome_exato"] = []string{*req.Nome}
	p.Filtros["usuario_id"] = []string{usuarioID.String()}
	p.Limite = 1
	lista, err := repo.ListarCartoes(&p)
	if err != nil {
		return utils.Wrap(err, msgErrPadraoListagem)
	}

	if len(lista.Dados) > 0 {
		if lista.Dados[0].DataDesativacao != nil {
			return utils.NewErr("Cartão desativado!")
		} else if *id == *lista.Dados[0].ID {
			return
		} else {
			return utils.NewErr("Já existe um cartão ativo com esse nome!")
		}
	}

	if err = repo.AtualizarCartao(reqInfra, id, usuarioID); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}

// RemoverCartao contém a regra de negócio para desativar um cartão de crédito
func RemoverCartao(id, usuarioID *uuid.UUID) (err error) {
	const msgErrPadrao = "Erro ao desativar cartão"

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	defer db.Close()

	repo := cartao.NovoRepo(db)

	if err := repo.RemoverCartao(id, usuarioID); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}

// ReativarCartao contém a regra de negócio para reativar um cartão
func ReativarCartao(id, usuarioID *uuid.UUID) (err error) {
	const msgErrPadrao = "Erro ao reativar cartão"

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	defer db.Close()

	repo := cartao.NovoRepo(db)

	if err := repo.ReativarCartao(id, usuarioID); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}
