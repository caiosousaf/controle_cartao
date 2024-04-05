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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// removerCategoria
func removerCategoria(c *gin.Context) {
	idCategoria, err := utils.GetUUIDFromParam(c, "categoria_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := categorias.RemoverCategoria(idCategoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// reativarCategoria
func reativarCategoria(c *gin.Context) {
	idCategoria, err := utils.GetUUIDFromParam(c, "categoria_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := categorias.ReativarCategoria(idCategoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}
