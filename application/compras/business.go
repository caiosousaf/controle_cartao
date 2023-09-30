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
)

// CadastrarCompra contém a regra de negócio para cadastrar uma nova compra
func CadastrarCompra(req *Req, idFatura *uuid.UUID) (idCompra *uuid.UUID, err error) {
	const (
		msgErrPadrao          = "Erro ao cadastrar nova compra"
		msgErrProxFaturas     = "Erro ao obter as próximas faturas"
		msgErrVerificarFatura = "Erro ao verificar fatura"
		msgErrCadastrarFatura = "Erro ao cadastrar nova fatura"
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
	}

	return
}
