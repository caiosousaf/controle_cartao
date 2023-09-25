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
