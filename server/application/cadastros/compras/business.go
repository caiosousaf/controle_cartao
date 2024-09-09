package compras

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/compras"
	"controle_cartao/domain/cadastros/faturas"
	infra "controle_cartao/infrastructure/cadastros/compras"
	infraFaturas "controle_cartao/infrastructure/cadastros/faturas"
	"controle_cartao/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// CadastrarCompra contém a regra de negócio para cadastrar uma nova compra
func CadastrarCompra(req *Req, idFatura, usuarioID *uuid.UUID) (idCompra *uuid.UUID, err error) {
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

	buscaFatura, err := repoFatura.BuscarFatura(idFatura, usuarioID)
	if err != nil {
		return idCompra, utils.Wrap(err, msgErrBuscarFatura)
	}

	if ok := repoFatura.CartaoPertenceUsuario(buscaFatura.FaturaCartaoID, usuarioID); ok != true {
		return idCompra, utils.NewErr("Cartão selecionado não pertence ao usuário")
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
		faturaID, err := repoFatura.VerificarFaturaCartao(&datas[i], idCartao, usuarioID)
		if errors.Is(err, sql.ErrNoRows) && faturaID == nil {
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
func ListarCompras(params *utils.Parametros, usuarioID *uuid.UUID) (res *ResComprasPag, err error) {
	const msgErrPadrao = "Erro ao listar compras"

	res = new(ResComprasPag)

	db, err := database.Conectar()
	if err != nil {
		return res, err
	}
	defer db.Close()

	repo := compras.NovoRepo(db)

	listaCompras, err := repo.ListarCompras(params, usuarioID)
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
func ObterTotalComprasValor(params *utils.Parametros, usuarioID *uuid.UUID) (res *ResTotalComprasValor, err error) {
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

	totalCompras, err := repo.ObterTotalComprasValor(params, usuarioID)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	*totalCompras.Total = "R$ " + *totalCompras.Total

	res = &ResTotalComprasValor{Total: totalCompras.Total}

	return
}

func PdfComprasFaturaCartao(params *utils.Parametros, usuarioID *uuid.UUID) (pdf *gofpdf.Fpdf, err error) {
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
		repo             = compras.NovoRepo(db)
		repoFatura       = faturas.NovoRepo(db)
		paramsFatura     = new(utils.Parametros)
		valorTotalFatura float64
	)

	params.Limite = utils.MaxLimit
	paramsFatura.Limite = utils.MaxLimit

	if !params.TemFiltro("fatura_id") && params.TemFiltro("cartao_id") {
		if params.TemFiltro("ano_exato") {
			paramsFatura.AdicionarFiltro("ano_exato", params.Filtros["ano_exato"][0])
		}

		cartaoIdStr := params.Filtros["cartao_id"][0]
		cartaoUuid, erro := uuid.Parse(cartaoIdStr)
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		listaFaturas, erro := repoFatura.ListarFaturasCartao(paramsFatura, &cartaoUuid, usuarioID)
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		totalFatura := make([]string, len(listaFaturas.Dados))
		var total float64

		pdf, erro = gerarPdf()
		if erro != nil {
			return pdf, utils.Wrap(erro, msgErrPadrao)
		}

		for j := range listaFaturas.Dados {
			params.AdicionarFiltro("fatura_id", listaFaturas.Dados[j].ID.String())
			valor, err := repo.ObterTotalComprasValor(params, usuarioID)
			if err != nil {
				return pdf, utils.Wrap(err, msgErrPadrao)
			}

			totalFatura[j] = *valor.Total
			totalFloat, err := strconv.ParseFloat(*valor.Total, 64)
			if err != nil {
				return pdf, utils.Wrap(err, "Não foi possível converter valor")
			}

			total += totalFloat

			params.RemoverFiltros("fatura_id")
		}

		tabelaMesesFaturasCartao(pdf, listaFaturas, tamanhoCel, alturaCel, totalFatura, total)

		for i := range listaFaturas.Dados {
			params.AdicionarFiltro("fatura_id", listaFaturas.Dados[i].ID.String())

			listaCompras, err := repo.ListarCompras(params, usuarioID)
			if err != nil {
				return pdf, utils.Wrap(err, msgErrPadrao)
			}

			tabelaFaturasPdf(pdf, listaCompras, nil, tamanhoCel, alturaCel)

			params.RemoverFiltros("fatura_id")
		}

		return
	}

	listaCompras, err := repo.ListarCompras(params, usuarioID)
	if err != nil {
		return pdf, utils.Wrap(err, msgErrPadrao)
	}

	for _, compra := range listaCompras.Dados {
		valorTotalFatura += *compra.ValorParcela
	}

	valorTotalFaturaString := strconv.FormatFloat(valorTotalFatura, 'f', -1, 64)

	pdf, err = gerarPdf()

	log.Println(pdf)

	tabelaFaturasPdf(pdf, listaCompras, &valorTotalFaturaString, tamanhoCel, alturaCel)

	return
}

func gerarPdf() (pdf *gofpdf.Fpdf, err error) {
	pdf = gofpdf.New("P", "mm", "A4", "")

	fontPaths := []string{
		"/app/font/CaviarDreams.ttf",
		"/app/font/CaviarDreams_Bold.ttf",
		"/app/font/CaviarDreams.ttf",
		"/app/font/CaviarDreams_Bold.ttf",
		"app/font/CaviarDreams.ttf",
		"app/font/CaviarDreams_Bold.ttf",
		"/server/font/CaviarDreams.ttf",
		"/server/font/CaviarDreams_Bold.ttf",
		"server/font/CaviarDreams.ttf",
		"server/font/CaviarDreams_Bold.ttf",
		"font/CaviarDreams.ttf",
		"font/CaviarDreams_Bold.ttf",
		"/font/CaviarDreams.ttf",
		"/font/CaviarDreams_Bold.ttf",
	}

	for _, path := range fontPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("Font file not found: %s", path)
		} else {
			log.Printf("Font file found: %s", path)
		}
	}

	pdf.AddUTF8Font("Caviar", "", "//app/font/CaviarDreams.ttf")
	pdf.AddUTF8Font("Caviar Bold", "B", "//app/font/CaviarDreams_Bold.ttf")
	pdf.AddUTF8Font("Caviar Italic", "I", "//app/font/CaviarDreams_Italic.ttf")
	pdf.AddUTF8Font("Caviar BoldItalic", "BI", "//app/font/CaviarDreams_BoldItalic.ttf")

	// Configura a fonte
	pdf.SetFont("Caviar", "", 5)
	pdf.AddPage()

	return
}

// tabelaFaturasPdf é a função responsável por montar o pdf com as faturas e suas compras
func tabelaFaturasPdf(pdf *gofpdf.Fpdf, listaCompras *infra.ComprasPag, valorTotalFatura *string, tamanho, altura float64) {
	left := float64(40)

	leftTitulo := float64(90)
	pdf.SetX(leftTitulo)

	pdf.SetFont("Caviar Bold", "B", 12)

	if valorTotalFatura != nil {
		pdf.CellFormat(tamanho, altura, fmt.Sprintf("Fatura do mês de %s - Valor total: R$ %s", *listaCompras.Dados[0].NomeFatura, *valorTotalFatura), "0", 0, "C", false, 0, "")
	} else {
		pdf.CellFormat(tamanho, altura, fmt.Sprintf("Fatura do mês de %s", *listaCompras.Dados[0].NomeFatura), "0", 0, "C", false, 0, "")
	}

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetFillColor(68, 68, 68)
	pdf.SetTextColor(255, 255, 255)

	pdf.SetX(left)

	for _, str := range colunasFaturasPdf {
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

	pdf.SetFont("Caviar", "", 5)

	return
}

// tabelaMesesFaturasCartao
func tabelaMesesFaturasCartao(pdf *gofpdf.Fpdf, listaFaturas *infraFaturas.FaturaPag, tamanho, altura float64, totalFatura []string, total float64) {
	var (
		leftMeses = float64(71)
		left      = float64(83)
	)

	pdf.SetFont("Caviar Bold", "B", 12)
	pdf.SetX(left)

	pdf.CellFormat(tamanho, altura, fmt.Sprintf("Cartão %s", *listaFaturas.Dados[0].NomeCartao), "0", 0, "C", false, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetFillColor(68, 68, 68)
	pdf.SetTextColor(255, 255, 255)

	pdf.SetX(leftMeses)
	for _, colunas := range colunasMesesFaturasCartao {
		pdf.SetFont("Caviar", "", 5)
		pdf.CellFormat(tamanho, altura, colunas, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	for i, fatura := range listaFaturas.Dados {
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(leftMeses)

		pdf.SetFont("Caviar", "", 5)

		pdf.CellFormat(tamanho, altura, *fatura.Nome, "1", 0, "C", false, 0, "")
		pdf.CellFormat(tamanho, altura, *fatura.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(tamanho, altura, totalFatura[i], "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.SetX(leftMeses)
	pdf.CellFormat(tamanho*2, altura, "Valor total", "1", 0, "C", true, 0, "")
	pdf.CellFormat(tamanho, altura, strconv.FormatFloat(total, 'f', -1, 64), "1", 0, "C", true, 0, "")

	pdf.Ln(-1)
	pdf.Ln(-1)
}
