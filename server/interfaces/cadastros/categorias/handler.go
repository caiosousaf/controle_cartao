package categorias

import (
	"controle_cartao/application/cadastros/categorias"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// listarCategorias
func listarCategorias(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	res, err := categorias.ListarCategorias(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
