package recorrente

import (
	"controle_cartao/application/cadastros/compras/recorrente"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// listarComprasRecorrentes godoc
func listarComprasRecorrentes(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := recorrente.ListarComprasRecorrentes(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// cadastrarComprasRecorrentes godoc
func cadastrarComprasRecorrentes(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	if err := recorrente.CadastrarComprasRecorrentes(usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// obterPrevisaoGastos godoc
func obterPrevisaoGastos(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	res, err := recorrente.ObterPrevisaoGastos(usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
