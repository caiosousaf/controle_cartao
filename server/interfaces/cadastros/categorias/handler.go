package categorias

import (
	"controle_cartao/application/cadastros/categorias"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// cadastrarCategoria
func cadastrarCategoria(c *gin.Context) {
	var req categorias.ReqCategoria
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	req.UsuarioID = middlewares.AuthUsuario(c)

	id, err := categorias.CadastrarCategoria(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, id)
}

// atualizarCategoria
func atualizarCategoria(c *gin.Context) {
	idCategoria, err := utils.GetUUIDFromParam(c, "categoria_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	var req categorias.ReqCategoria
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	req.UsuarioID = middlewares.AuthUsuario(c)

	if err = categorias.AtualizarCategoria(&req, idCategoria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()

	}
}

// listarCategorias
func listarCategorias(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := categorias.ListarCategorias(&params, usuarioID)
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

	usuarioID := middlewares.AuthUsuario(c)

	if err := categorias.RemoverCategoria(idCategoria, usuarioID); err != nil {
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

	usuarioID := middlewares.AuthUsuario(c)

	if err := categorias.ReativarCategoria(idCategoria, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// buscarCategoria
func buscarCategoria(c *gin.Context) {
	idCategoria, err := utils.GetUUIDFromParam(c, "categoria_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := categorias.BuscarCategoria(idCategoria, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
