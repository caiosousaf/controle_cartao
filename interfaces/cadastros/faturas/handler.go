package faturas

import (
	"controle_cartao/application/cadastros/faturas"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// listarFaturasCartao godoc
func listarFaturasCartao(c *gin.Context) {
	p, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	id, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	res, err := faturas.ListarFaturasCartao(&p, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// buscarFaturaCartao godoc
func buscarFaturaCartao(c *gin.Context) {
	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idCartao, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	res, err := faturas.BuscarFaturaCartao(idFatura, idCartao)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// cadastrarFatura godoc
func cadastrarFatura(c *gin.Context) {
	var req faturas.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idCartao, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	req.FaturaCartaoID = idCartao

	id, err := faturas.CadastrarFatura(&req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, id)
}
