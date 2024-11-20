package cartao

import (
	"controle_cartao/application/cadastros/cartao"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// cadastrarCartao godoc
func cadastrarCartao(c *gin.Context) {
	var req cartao.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	req.UsuarioID = middlewares.AuthUsuario(c)

	id, err := cartao.CadastrarCartao(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, id)
}

// listarCartoes godoc
func listarCartoes(c *gin.Context) {
	p, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := cartao.ListarCartoes(&p, usuarioID)
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

	usuarioID := middlewares.AuthUsuario(c)

	res, err := cartao.BuscarCartao(id, usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// atualizarCartao godoc
func atualizarCartao(c *gin.Context) {
	var req cartao.ReqAtualizar
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	id, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	if err := cartao.AtualizarCartao(&req, id, usuarioID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// removerCartao godoc
func removerCartao(c *gin.Context) {
	id, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	if err := cartao.RemoverCartao(id, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// reativarCartao godoc
func reativarCartao(c *gin.Context) {
	id, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	if err := cartao.ReativarCartao(id, usuarioID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
