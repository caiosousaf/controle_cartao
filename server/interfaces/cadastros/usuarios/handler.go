package usuarios

import (
	"controle_cartao/application/usuarios"
	"github.com/gin-gonic/gin"
	"net/http"
)

// cadastrarUsuario
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

// login
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
