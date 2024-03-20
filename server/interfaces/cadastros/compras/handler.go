package compras

import (
	"controle_cartao/application/cadastros/compras"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// cadastrarCompra godoc
func cadastrarCompra(c *gin.Context) {
	var req compras.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idCompra, err := compras.CadastrarCompra(&req, idFatura)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, idCompra)
}

// listarCompras godoc
func listarCompras(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	res, err := compras.ListarCompras(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// pdfComprasFaturaCartao godoc
func pdfComprasFaturaCartao(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	pdf, err := compras.PdfComprasFaturaCartao(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	// Define o cabeçalho Content-Disposition para indicar que o arquivo é um anexo
	c.Header("Content-Disposition", "attachment; filename=compras.pdf")

	// Define o tipo de conteúdo como PDF
	c.Header("Content-Type", "application/pdf")

	// Gera o PDF e escreve no contexto de resposta
	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, "Erro ao gerar PDF")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "PDF gerado com sucesso")
}

// obterTotalComprasValor godoc
func obterTotalComprasValor(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	res, err := compras.ObterTotalComprasValor(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
