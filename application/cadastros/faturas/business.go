package faturas

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/faturas"
	"controle_cartao/utils"
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
func BuscarFaturaCartao(idFatura, idCartao *uuid.UUID) (res *Res, err error) {
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
