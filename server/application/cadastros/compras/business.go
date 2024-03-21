package compras

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/compras"
	"controle_cartao/domain/cadastros/faturas"
	infra "controle_cartao/infrastructure/cadastros/compras"
	infraFaturas "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
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

	if params.TemFiltro("ultima_parcela") && !params.TemFiltro("data_especifica") {
		err = utils.NewErr("Erro ao filtrar, filtro 'data_especifica' é obrigatório ao ser passado o filtro 'ultima_parcela'")
		return res, utils.Wrap(err, msgErrPadrao)
	}

	totalCompras, err := repo.ObterTotalComprasValor(params)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	*totalCompras.Total = "R$ " + *totalCompras.Total

	res = &ResTotalComprasValor{Total: totalCompras.Total}

	return
}

func PdfComprasFaturaCartao(params *utils.Parametros) (pdf *gofpdf.Fpdf, err error) {
	const (
		tamanhoCel   = float64(18)
		alturaCel    = float64(5)
		msgErrPadrao = "Erro ao gerar pdf"
	)

	db, err := database.Conectar()
	if err != nil {
		return pdf, err
	}
	defer db.Close()

	var (
		repo         = compras.NovoRepo(db)
		repoFatura   = faturas.NovoRepo(db)
		paramsFatura = new(utils.Parametros)
	)

	params.Limite = utils.MaxLimit
	paramsFatura.Limite = utils.MaxLimit

	if !params.TemFiltro("fatura_id") && params.TemFiltro("cartao_id") {
		cartaoIdStr := params.Filtros["cartao_id"][0]
		cartaoUuid, erro := uuid.Parse(cartaoIdStr)
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		listaFaturas, erro := repoFatura.ListarFaturasCartao(paramsFatura, &cartaoUuid)
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		pdf, erro = gerarPdf()
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		for i := range listaFaturas.Dados {
			params.AdicionarFiltro("fatura_id", listaFaturas.Dados[i].ID.String())

			listaCompras, err := repo.ListarCompras(params)
			if err != nil {
				return pdf, utils.Wrap(err, msgErrPadrao)
			}

			basePdf(pdf, listaCompras, tamanhoCel, alturaCel, colunasPdf)

			params.RemoverFiltros("fatura_id")
		}

		return
	}

	listaCompras, err := repo.ListarCompras(params)
	if err != nil {
		return pdf, utils.Wrap(err, msgErrPadrao)
	}

	pdf, err = gerarPdf()

	basePdf(pdf, listaCompras, tamanhoCel, alturaCel, colunasPdf)

	return
}

func gerarPdf() (pdf *gofpdf.Fpdf, err error) {
	pdf = gofpdf.New("P", "mm", "A4", "")

	pdf.AddUTF8Font("Caviar", "", "server/font/CaviarDreams.ttf")
	pdf.AddUTF8Font("Caviar Bold", "B", "server/font/CaviarDreams_Bold.ttf")
	pdf.AddUTF8Font("Caviar Italic", "I", "server/font/CaviarDreams_Italic.ttf")
	pdf.AddUTF8Font("Caviar BoldItalic", "BI", "server/font/CaviarDreams_BoldItalic.ttf")

	// Configura a fonte
	pdf.SetFont("Caviar", "", 5)
	pdf.AddPage()

	return
}

func basePdf(pdf *gofpdf.Fpdf, listaCompras *infra.ComprasPag, tamanho, altura float64, header []string) {
	fancyTable := func() {

		left := float64(40)

		leftTitulo := float64(90)
		pdf.SetX(leftTitulo)

		pdf.SetFont("Caviar Bold", "B", 12)

		pdf.CellFormat(tamanho, altura, fmt.Sprintf("Fatura do mês de %s", *listaCompras.Dados[0].NomeFatura), "0", 0, "C", false, 0, "")

		pdf.Ln(-1)
		pdf.Ln(-1)

		pdf.SetFillColor(68, 68, 68)
		pdf.SetTextColor(255, 255, 255)
		// Cor das linhas
		//pdf.SetDrawColor(128, 0, 0)
		//pdf.SetLineWidth(.3)

		pdf.SetX(left)

		for _, str := range header {
			pdf.SetFont("Caviar", "", 5)
			pdf.CellFormat(tamanho, altura, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)

		for _, compra := range listaCompras.Dados {
			pdf.SetX(left)

			dataCompraFormat := *compra.DataCompra

			dataCompraFormat = dataCompraFormat[:10]

			pdf.SetFont("Caviar", "", 5)

			pdf.CellFormat(tamanho, altura, *compra.Nome, "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, *compra.LocalCompra, "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, *compra.CategoriaNome, "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, strconv.FormatFloat(*compra.ValorParcela, 'f', -1, 64), "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, strconv.FormatInt(*compra.ParcelaAtual, 10), "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, strconv.FormatInt(*compra.QuantidadeParcelas, 10), "1", 0, "C", false, 0, "")
			pdf.CellFormat(tamanho, altura, dataCompraFormat, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		pdf.Ln(-1)
	}

	pdf.SetFont("Caviar", "", 5)
	fancyTable()

	return
}
