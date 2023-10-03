package compras

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/compras"
	"controle_cartao/domain/cadastros/faturas"
	infra "controle_cartao/infrastructure/cadastros/compras"
	infraFaturas "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

// CadastrarCompra contém a regra de negócio para cadastrar uma nova compra
func CadastrarCompra(req *Req, idFatura *uuid.UUID) (idCompra *uuid.UUID, err error) {
	const (
		msgErrPadrao                  = "Erro ao cadastrar nova compra"
		msgErrProxFaturas             = "Erro ao obter as próximas faturas"
		msgErrVerificarFatura         = "Erro ao verificar fatura"
		msgErrCadastrarFatura         = "Erro ao cadastrar nova fatura"
		msgErrBuscarFatura            = "Erro ao buscar fatura"
		formatoDataEsperado           = "2006-01-02"
		formatoDataVencimentoEsperado = "2006-01-02T15:04:05Z07:00"
	)

	var (
		reqInfra       = new(infra.Compras)
		reqInfraFatura = new(infraFaturas.Fatura)
	)

	db, err := database.Conectar()
	if err != nil {
		return idCompra, err
	}
	defer db.Close()

	var (
		repo       = compras.NovoRepo(db)
		repoFatura = faturas.NovoRepo(db)
	)

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return idCompra, utils.Wrap(err, msgErrPadrao)
	}

	buscaFatura, err := repoFatura.BuscarFatura(idFatura)
	if err != nil {
		return idCompra, utils.Wrap(err, msgErrBuscarFatura)
	}

	dataCompraData, err := time.Parse(formatoDataEsperado, *reqInfra.DataCompra)
	if err != nil {
		return nil, utils.NewErr("Erro ao converter data compra. A data inserida é inválida. Ex: 2006-01-02")
	}

	dataVencimentoData, err := time.Parse(formatoDataVencimentoEsperado, *buscaFatura.DataVencimento)
	if err != nil {
		return idCompra, utils.NewErr("Erro ao converter data vencimento. A data inserida é inválida. Ex: 2006-01-02T15:04:05Z07:00")
	}

	if dataVencimentoData.Before(dataCompraData) {
		return idCompra, utils.NewErr("Data da compra deve ser menor que a data de vencimento da fatura")
	}

	datas, meses, idCartao, err := repoFatura.ObterProximasFaturas(req.ParcelaAtual, req.QuantidadeParcelas, idFatura)
	if err != nil {
		return idCompra, utils.Wrap(err, msgErrProxFaturas)
	}

	for i := range datas {
		faturaID, err := repoFatura.VerificarFaturaCartao(&datas[i], idCartao)
		if err == sql.ErrNoRows && faturaID == nil {
			reqInfraFatura.Nome = &meses[i]
			reqInfraFatura.DataVencimento = &datas[i]
			reqInfraFatura.FaturaCartaoID = idCartao

			if err := repoFatura.CadastrarFatura(reqInfraFatura); err != nil {
				return idCompra, utils.Wrap(err, msgErrCadastrarFatura)
			}

			if reqInfraFatura.ID != nil {
				idFatura = reqInfraFatura.ID
			}

		} else if faturaID != nil {
			idFatura = faturaID
		} else {
			return idCompra, utils.Wrap(err, msgErrVerificarFatura)
		}

		reqInfra.FaturaID = idFatura
		if err := repo.CadastrarCompra(reqInfra); err != nil {
			return idCompra, utils.Wrap(err, msgErrPadrao)
		}

		*req.ParcelaAtual++

		idCompra = reqInfra.ID
	}

	return
}

// ListarCompras contém a regra de negócio para listar as compras
func ListarCompras(params *utils.Parametros) (res *ResComprasPag, err error) {
	const msgErrPadrao = "Erro ao listar compras"

	res = new(ResComprasPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := compras.NovoRepo(db)

	listaCompras, err := repo.ListarCompras(params)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res.Dados = make([]ResCompras, len(listaCompras.Dados))
	for i := 0; i < len(listaCompras.Dados); i++ {
		if err = utils.ConvertStructByAlias(&listaCompras.Dados[i], &res.Dados[i]); err != nil {
			return res, utils.Wrap(err, msgErrPadrao)
		}
	}

	res.Total, res.Prox = listaCompras.Total, listaCompras.Prox

	return
}

// ObterTotalComprasValor contém a regra de negócio para obter o total das compras
func ObterTotalComprasValor(params *utils.Parametros) (res *ResTotalComprasValor, err error) {
	const msgErrPadrao = "Erro ao obter o total das compras"

	res = new(ResTotalComprasValor)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := compras.NovoRepo(db)

	totalCompras, err := repo.ObterTotalComprasValor(params)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	utils.ArredondarParaDuasCasasDecimais(totalCompras.Total)

	res = &ResTotalComprasValor{Total: totalCompras.Total}

	return
}
