package cartao

import (
	"controle_cartao/application/cadastros/cartao"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// listarCartoes godoc
func listarCartoes(c *gin.Context) {
	p, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	res, err := cartao.ListarCartoes(&p)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, res)
}

// buscarCartao godoc
func buscarCartao(c *gin.Context) {
	id, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	res, err := cartao.BuscarCartao(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
