package recorrente

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/compras/recorrente"
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
