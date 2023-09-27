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

// BuscarFaturaCartao contém a regra de negócio para buscar uma fatura de um cartão dado os id's fornecidos
func BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (res *Res, err error) {
	const msgErrPadrao = "Erro ao buscar fatura de cartão"

	res = new(Res)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := faturas.NovoRepo(db)

	buscaFatura, err := repo.BuscarFaturaCartao(idFatura, idCartao)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res = &Res{
		ID:             buscaFatura.ID,
		Nome:           buscaFatura.Nome,
		FaturaCartaoID: buscaFatura.FaturaCartaoID,
		NomeCartao:     buscaFatura.NomeCartao,
		DataCriacao:    buscaFatura.DataCriacao,
		DataVencimento: buscaFatura.DataVencimento,
	}

	return
}

// CadastrarFatura contém a regra de negócio para cadastrar uma nova fatura
func CadastrarFatura(req *Req) (id *uuid.UUID, err error) {
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
