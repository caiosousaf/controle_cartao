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
