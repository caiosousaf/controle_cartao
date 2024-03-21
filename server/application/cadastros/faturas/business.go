package faturas

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/faturas"
	infra "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
)

// ListarFaturasCartao contém a regra de negócio para listar as faturas de um cartão
func ListarFaturasCartao(p *utils.Parametros, id *uuid.UUID) (res *ResPag, err error) {
	const msgErrPadrao = "Erro ao listar faturas de um cartão"

	res = new(ResPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	listaFaturas, err := repo.ListarFaturasCartao(p, id)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res.Dados = make([]Res, len(listaFaturas.Dados))
	for i := 0; i < len(listaFaturas.Dados); i++ {
		if err = utils.ConvertStructByAlias(&listaFaturas.Dados[i], &res.Dados[i]); err != nil {
			return res, utils.Wrap(err, msgErrPadrao)
		}
	}

	res.Total, res.Prox = listaFaturas.Total, listaFaturas.Prox

	return
}

// BuscarFatura contém a regra de negócio para buscar uma fatura
func BuscarFatura(idFatura *uuid.UUID) (res *Res, err error) {
	const msgErrPadrao = "Erro ao buscar fatura de cartão"

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	buscaFatura, err := repo.BuscarFatura(idFatura)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res = &Res{
		ID:             buscaFatura.ID,
		Nome:           buscaFatura.Nome,
		FaturaCartaoID: buscaFatura.FaturaCartaoID,
		NomeCartao:     buscaFatura.NomeCartao,
		Status:         buscaFatura.Status,
		DataCriacao:    buscaFatura.DataCriacao,
		DataVencimento: buscaFatura.DataVencimento,
	}

	return
}

// CadastrarFatura contém a regra de negócio para cadastrar uma nova fatura
func CadastrarFatura(req *Req, idCartao *uuid.UUID) (id *uuid.UUID, err error) {
	const (
		msgErrPadrao                = "Erro ao cadastrar nova fatura"
		msgErrPadraoVerificarFatura = "Erro ao verificar se já existe fatura cadastrada para o mês de vencimento"
	)

	var reqInfra = new(infra.Fatura)

	db, err := database.Conectar()
	if err != nil {
		return id, err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	numeroMesVencimento, err := utils.ObterNumeroDoMes(*req.DataVencimento)
	if err != nil {
		return nil, err
	}

	mesStringFormatado := utils.NumeroParaNomeMes(numeroMesVencimento)

	req.Nome = utils.GetStringPointer(mesStringFormatado)
	req.FaturaCartaoID = idCartao

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return id, utils.Wrap(err, msgErrPadrao)
	}

	cartaoID, err := repo.VerificarFaturaCartao(req.DataVencimento, req.FaturaCartaoID)

	if err != nil {
		if err == sql.ErrNoRows && cartaoID == nil {
			if err = repo.CadastrarFatura(reqInfra); err != nil {
				return id, utils.Wrap(err, msgErrPadrao)
			}

			id = reqInfra.ID
		} else {
			return id, utils.Wrap(err, msgErrPadraoVerificarFatura)
		}
	} else {
		return id, utils.NewErr("Fatura do mês escolhido já existe")
	}

	return
}

// AtualizarFatura contém a regra de negócio para atualizar uma fatura
func AtualizarFatura(req *ReqAtualizar, idCartao, idFatura *uuid.UUID) (err error) {
	const (
		msgErrPadrao                = "Erro ao atualizar fatura"
		msgErrPadraoVerificarFatura = "Erro ao verificar se já existe fatura cadastrada para o mês de vencimento"
	)

	var reqInfra = new(infra.Fatura)

	db, err := database.Conectar()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	numeroMesVencimento, err := utils.ObterNumeroDoMes(*req.DataVencimento)
	if err != nil {
		return err
	}

	mesStringFormatado := utils.NumeroParaNomeMes(numeroMesVencimento)

	req.Nome = utils.GetStringPointer(mesStringFormatado)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	faturaID, err := repo.VerificarFaturaCartao(req.DataVencimento, idCartao)

	if err != nil {
		if err == sql.ErrNoRows && faturaID == nil {
			if err = repo.AtualizarFatura(reqInfra, idFatura); err != nil {
				return utils.Wrap(err, msgErrPadrao)
			}

		} else {
			return utils.Wrap(err, msgErrPadraoVerificarFatura)
		}
	} else {
		return utils.NewErr("Fatura do mês escolhido já existe")
	}

	return
}

// AtualizarStatusFatura contém a regra de negócio para atualizar o status de uma fatura
func AtualizarStatusFatura(req *ReqAtualizarStatus, idFatura *uuid.UUID) (err error) {
	const (
		msgErrPadrao       = "Erro ao atualizar status de fatura"
		msgErrBuscarFatura = "Erro ao buscar fatura"
	)

	var reqInfra = new(infra.Fatura)

	db, err := database.Conectar()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	buscaFatura, err := repo.BuscarFatura(idFatura)
	if err != nil {
		return utils.Wrap(err, msgErrBuscarFatura)
	}

	if *buscaFatura.Status != FaturaPaga {
		if *buscaFatura.Status == FaturaAtrasada && *req.Status != FaturaEmAberto {
			if err = repo.AtualizarStatusFatura(reqInfra, idFatura); err != nil {
				return utils.Wrap(err, msgErrPadrao)
			}
		} else if *buscaFatura.Status == FaturaAtrasada && *req.Status == FaturaEmAberto {
			return utils.NewErr("Operação inválida, Fatura se encontra em atraso")
		} else if *buscaFatura.Status == FaturaEmAberto {
			if err = repo.AtualizarStatusFatura(reqInfra, idFatura); err != nil {
				return utils.Wrap(err, msgErrPadrao)
			}
		}

	} else {
		return utils.NewErr("Não é possível alterar status, Fatura já paga")
	}

	if err = repo.AtualizarStatusFatura(reqInfra, idFatura); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}
