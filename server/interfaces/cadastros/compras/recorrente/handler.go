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

// cadastrarNovaCompraRecorrente
func cadastrarNovaCompraRecorrente(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	var req recorrente.Recorrentes

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := recorrente.CadastrarNovaCompraRecorrente(&req, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// atualizarCompraRecorrente
func atualizarCompraRecorrente(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	recorrenteID, err := utils.GetUUIDFromParam(c, "recorrente_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	var req recorrente.Recorrentes

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := recorrente.AtualizarCompraRecorrente(&req, recorrenteID, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// desativarCompraRecorrente godoc
func desativarCompraRecorrente(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	recorrenteID, err := utils.GetUUIDFromParam(c, "recorrente_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := recorrente.DesativarCompraRecorrente(recorrenteID, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// reativarCompraRecorrente godoc
func reativarCompraRecorrente(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	recorrenteID, err := utils.GetUUIDFromParam(c, "recorrente_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := recorrente.ReativarCompraRecorrente(recorrenteID, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// removerCompraRecorrente godoc
func removerCompraRecorrente(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	recorrenteID, err := utils.GetUUIDFromParam(c, "recorrente_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := recorrente.RemoverCompraRecorrente(recorrenteID, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
