package postgres

import (
	model "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

// DBFatura é uma estrutura base para acesso aos metodos do banco postgres para manipulação das faturas
type DBFatura struct {
	DB *sql.DB
}

// ListarFaturasCartao lista de forma paginada os dados das faturas de um cartão no banco de dados
func (pg *DBFatura) ListarFaturasCartao(p *utils.Parametros, id, usuarioID *uuid.UUID) (res *model.FaturaPag, err error) {
	var t model.Fatura

	res = new(model.FaturaPag)

	campos, _, err := p.ValidFields(&t)
	if err != nil {
		return res, err
	}

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(campos...).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC on TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.fatura_cartao_id": id,
			"TC.data_desativacao":  nil,
			"TC.usuario_id":        usuarioID,
		}).OrderBy("TFC.data_criacao")

	consultaComFiltro := p.CriarFiltros(consultaSql, map[string]utils.Filtro{
		"nome_exato": utils.CriarFiltros("LOWER(TF.nome) = LOWER(?)", utils.FlagFiltroEq),
		"ano_exato":  utils.CriarFiltros("EXTRACT(year from data_vencimento) = ?", utils.FlagFiltroEq),
		"pago":       utils.CriarFiltros("(TFC.status <> 'Pago') = ?", utils.FlagFiltroEq),
	}).PlaceholderFormat(sq.Dollar)

	dados, prox, total, err := utils.ConfigurarPaginacao(p, &t, &consultaComFiltro)
	if err != nil {
		return res, err
	}

	res.Dados, res.Prox, res.Total = dados.([]model.Fatura), prox, total

	return
}

// BuscarFaturaCartao busca os dados de uma fatura de um cartão no banco de dados dado os ID's fornecidos
func (pg *DBFatura) BuscarFaturaCartao(idCartao, idFatura *uuid.UUID) (res *model.Fatura, err error) {
	res = new(model.Fatura)
	if err = sq.StatementBuilder.RunWith(pg.DB).Select(`TFC.id, TFC.nome, TFC.fatura_cartao_id,
				TC.nome, TFC.status, TFC.data_criacao, TFC.data_vencimento`).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC on TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.fatura_cartao_id": idCartao,
			"TFC.id":               idFatura,
			"TC.data_desativacao":  nil,
		}).PlaceholderFormat(sq.Dollar).
		Scan(&res.ID, &res.Nome, &res.FaturaCartaoID, &res.NomeCartao, &res.Status, &res.DataCriacao, &res.DataVencimento); err != nil {
		return res, err
	}

	return
}

// BuscarFatura busca os dados de uma fatura de um cartão no banco de dados dado os ID's fornecidos
func (pg *DBFatura) BuscarFatura(idFatura, usuarioID *uuid.UUID) (res *model.Fatura, err error) {
	res = new(model.Fatura)
	if err = sq.StatementBuilder.RunWith(pg.DB).Select(`TFC.id, TFC.nome, TFC.fatura_cartao_id,
				TC.nome, TFC.status, TFC.data_criacao, TFC.data_vencimento`).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC on TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.id":        idFatura,
			"TC.usuario_id": usuarioID,
		}).PlaceholderFormat(sq.Dollar).
		Scan(&res.ID, &res.Nome, &res.FaturaCartaoID, &res.NomeCartao, &res.Status, &res.DataCriacao, &res.DataVencimento); err != nil {
		return res, err
	}

	return
}

// ObterProximasFaturas obtém as próximas possíveis datas de faturas de um cartão no banco de dados
func (pg *DBFatura) ObterProximasFaturas(parcela_atual, qtd_parcelas *int64, idFatura *uuid.UUID) (datas, meses []string, idCartao *uuid.UUID, err error) {
	var (
		data       time.Time
		mes        int
		valorLimit = *qtd_parcelas - *parcela_atual
	)
	consultaSelect := fmt.Sprintf(`GENERATE_SERIES(TFC.data_vencimento, TFC.data_vencimento + INTERVAL '%d months' - INTERVAL '1 month',
                		INTERVAL '1 month')::DATE, TC.id,
    EXTRACT(MONTH FROM GENERATE_SERIES(TFC.data_vencimento, TFC.data_vencimento + INTERVAL '%d months' - INTERVAL '1 month',
    INTERVAL '1 month')::DATE) AS numero_mes`, *qtd_parcelas, *qtd_parcelas)

	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select(consultaSelect).
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC ON TC.id = TFC.fatura_cartao_id").
		Where(sq.Eq{
			"TFC.id": idFatura,
		}).PlaceholderFormat(sq.Dollar).
		Limit(uint64(valorLimit) + 1)

	rows, err := consultaSql.Query()
	if err != nil {
		return datas, meses, idCartao, err
	}

	for rows.Next() {
		err := rows.Scan(&data, &idCartao, &mes)
		if err != nil {
			return datas, meses, idCartao, err
		}

		// Formata a data e transforma em string
		formattedDate := data.Format("2006-01-02")
		mesFormatado := utils.NumeroParaNomeMes(mes)

		datas = append(datas, formattedDate)
		meses = append(meses, mesFormatado)
	}

	return
}

// VerificarFaturaCartao verifica se existe fatura de um cartão para a data escolhida
func (pg *DBFatura) VerificarFaturaCartao(data *string, idCartao, usuarioID *uuid.UUID) (faturaID *uuid.UUID, err error) {
	consultaSql := sq.StatementBuilder.RunWith(pg.DB).Select("TFC.id").
		From("public.t_fatura_cartao TFC").
		Join("public.t_cartao TC ON TC.id = TFC.fatura_cartao_id").
		Where(fmt.Sprintf("EXTRACT(MONTH FROM TFC.data_vencimento) = EXTRACT(MONTH FROM '%s'::DATE)", *data)).
		Where(fmt.Sprintf("EXTRACT(YEAR FROM TFC.data_vencimento) = EXTRACT(YEAR FROM '%s'::DATE)", *data)).
		Where(sq.Eq{"TC.id": idCartao}).
		Where(sq.Eq{
			"TC.id":         idCartao,
			"TC.usuario_id": usuarioID,
		}).
		PlaceholderFormat(sq.Dollar)

	if err = consultaSql.QueryRow().Scan(&faturaID); err != nil {
		return faturaID, err
	}

	return
}

// CadastrarFatura cadastra uma nova fatura de cartão no banco de dados
func (pg *DBFatura) CadastrarFatura(req *model.Fatura) (err error) {
	if err = sq.StatementBuilder.RunWith(pg.DB).Insert("public.t_fatura_cartao").
		Columns("nome", "fatura_cartao_id", "data_vencimento").
		Values(req.Nome, req.FaturaCartaoID, req.DataVencimento).
		Suffix(`RETURNING "id"`).
		PlaceholderFormat(sq.Dollar).
		Scan(&req.ID); err != nil {
		return err
	}

	return
}

// AtualizarFatura atualiza uma fatura de cartão no banco de dados
func (pg *DBFatura) AtualizarFatura(req *model.Fatura, idFatura *uuid.UUID) (err error) {
	consultaUpdate := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_fatura_cartao").
		Set("nome", req.Nome).
		Set("data_vencimento", req.DataVencimento).
		Where(sq.Eq{
			"id": idFatura,
		}).
		PlaceholderFormat(sq.Dollar)

	result, err := consultaUpdate.Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("Fatura com o ID informado não existe")
	}

	return
}

// AtualizarStatusFatura atualiza o status de uma fatura de cartão no banco de dados
func (pg *DBFatura) AtualizarStatusFatura(req *model.Fatura, idFatura *uuid.UUID) (err error) {
	result, err := sq.StatementBuilder.RunWith(pg.DB).Update("public.t_fatura_cartao").
		Set("status", req.Status).
		Where(sq.Eq{
			"id": idFatura,
		}).
		PlaceholderFormat(sq.Dollar).
		Exec()
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.NewErr("Não foi possível identificar a fatura pelo ID informado")
	}

	return
}

// CartaoPertenceUsuario verifica se o cartão informado pertence ao usuário
func (pg *DBFatura) CartaoPertenceUsuario(idCartao, usuarioID *uuid.UUID) bool {
	var valor int64
	_ = sq.StatementBuilder.RunWith(pg.DB).Select(`1`).
		From("public.t_cartao TC").
		Where(sq.Eq{
			"TC.id":         idCartao,
			"TC.usuario_id": usuarioID,
		}).PlaceholderFormat(sq.Dollar).
		Scan(&valor)

	return valor > 0
}
