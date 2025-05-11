package recorrente

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/compras/recorrente"
	infra "controle_cartao/infrastructure/cadastros/compras/recorrente"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// ListarComprasRecorrentes contém a regra de negócio para listar as compras recorrentes
func ListarComprasRecorrentes(params *utils.Parametros, usuarioID *uuid.UUID) (res *ResRecorrentesPag, err error) {
	const msgErrPadrao = "Erro ao listar compras recorrentes"

	res = new(ResRecorrentesPag)

	db, err := database.Conectar()
	if err != nil {
		return nil, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := recorrente.NovoRepo(db)

	listaRecorrentes, err := repo.ListarComprasRecorrentes(params, usuarioID)
	if err != nil {
		return nil, utils.Wrap(err, msgErrPadrao)
	}

	res.Dados = make([]ResRecorrentes, len(listaRecorrentes.Dados))
	for i := range listaRecorrentes.Dados {
		if err = utils.ConvertStructByAlias(&listaRecorrentes.Dados[i], &res.Dados[i]); err != nil {
			return res, utils.Wrap(err, msgErrPadrao)
		}
	}

	res.Total, res.Prox = listaRecorrentes.Total, listaRecorrentes.Prox

	return
}

// CadastrarComprasRecorrentes contém a regra para cadastrar as compras recorrentes
func CadastrarComprasRecorrentes(usuarioID *uuid.UUID) (err error) {
	const msgErrPadrao = "Erro ao cadastrar compras recorrentes"

	var params = new(utils.Parametros)

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	defer tx.Rollback()

	repo := recorrente.NovoRepo(db)

	params.Limite = utils.MaxLimit
	comprasRecorrentes, err := repo.ListarComprasRecorrentes(params, usuarioID)
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	fatura, err := repo.ObterFaturaCartaoGeral(usuarioID)
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	for i := range comprasRecorrentes.Dados {
		if err := repo.CadastrarCompraRecorrente(&infra.ComprasRecorrentes{
			Nome:               comprasRecorrentes.Dados[i].Nome,
			Descricao:          comprasRecorrentes.Dados[i].Descricao,
			LocalCompra:        comprasRecorrentes.Dados[i].LocalCompra,
			CategoriaID:        comprasRecorrentes.Dados[i].CategoriaID,
			ValorParcela:       comprasRecorrentes.Dados[i].ValorParcela,
			ParcelaAtual:       utils.GetPointer[int64](1),
			QuantidadeParcelas: utils.GetPointer[int64](0),
			FaturaID:           fatura,
			DataCompra:         utils.GetPointer(comprasRecorrentes.Dados[i].DataCriacao.Format("2006-01-02")),
			Recorrente:         true,
		}); err != nil {
			return utils.Wrap(err, msgErrPadrao)
		}
	}

	if err = tx.Commit(); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}
