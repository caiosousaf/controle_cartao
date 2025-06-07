package usuarios

import (
	"controle_cartao/application/usuarios"
	"controle_cartao/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

// cadastrarUsuario godoc
func cadastrarUsuario(c *gin.Context) {
	var req usuarios.ReqUsuario
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res, err := usuarios.CadastrarUsuario(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, res)
}

// login godoc
func login(c *gin.Context) {
	var req usuarios.ReqUsuarioLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res, err := usuarios.LoginUsuario(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// atualizarSenhaUsuario godoc
func atualizarSenhaUsuario(c *gin.Context) {
	usuarioID := middlewares.AuthUsuario(c)

	var req usuarios.ReqAlterarSenhaUsuario

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	if err := usuarios.AtualizarSenhaUsuario(&req, usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
