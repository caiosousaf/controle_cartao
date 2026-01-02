package compras

import (
	"controle_cartao/application/cadastros/compras"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// cadastrarCompra godoc
func cadastrarCompra(c *gin.Context) {
	var req compras.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	idCompra, err := compras.CadastrarCompra(&req, idFatura, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, idCompra)
}

// antecipacaoParcelas godoc
func antecipacaoParcelas(c *gin.Context) {
	var req compras.ReqAntecipacaoParcelas
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	if err := compras.AntecipacaoParcelas(&req, idFatura, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// obterParcelasDisponiveisAntecipacao
func obterParcelasDisponiveisAntecipacao(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	identificadorCompra, err := utils.GetUUIDFromParam(c, "identificador_compra")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	faturaID, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	res, err := compras.ObterParcelasDisponiveisAntecipacao(identificadorCompra, faturaID, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// listarCompras godoc
func listarCompras(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := compras.ListarCompras(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// atualizarCompras
func atualizarCompras(c *gin.Context) {
	var req compras.Req

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	compraID, err := utils.GetUUIDFromParam(c, "compra_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	atualizarTodasParcelas, err := strconv.ParseBool(c.Param("atualizar_todas_parcelas"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := compras.AtualizarCompras(&req, usuarioID, compraID, atualizarTodasParcelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// removerCompra
func removerCompra(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	compraID, err := utils.GetUUIDFromParam(c, "compra_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	removerTodasParcelas, err := strconv.ParseBool(c.Param("remover_todas_parcelas"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := compras.RemoverCompra(compraID, usuarioID, removerTodasParcelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// pdfComprasFaturaCartao godoc
func pdfComprasFaturaCartao(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	pdf, err := compras.PdfComprasFaturaCartao(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	// Define o cabeçalho Content-Disposition para indicar que o arquivo é um anexo
	c.Header("Content-Disposition", "attachment; filename=compras.pdf")

	// Define o tipo de conteúdo como PDF
	c.Header("Content-Type", "application/pdf")

	// Gera o PDF e escreve no contexto de resposta
	if err = pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}
}

// obterTotalComprasValor godoc
func obterTotalComprasValor(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := compras.ObterTotalComprasValor(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
