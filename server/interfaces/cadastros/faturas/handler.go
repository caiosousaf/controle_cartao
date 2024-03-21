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

// buscarFatura godoc
func buscarFatura(c *gin.Context) {
	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	res, err := faturas.BuscarFatura(idFatura)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// cadastrarFatura godoc
func cadastrarFatura(c *gin.Context) {
	var req faturas.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idCartao, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	id, err := faturas.CadastrarFatura(&req, idCartao)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, id)
}

// atualizarFatura godoc
func atualizarFatura(c *gin.Context) {
	var req faturas.ReqAtualizar
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idCartao, err := utils.GetUUIDFromParam(c, "cartao_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := faturas.AtualizarFatura(&req, idCartao, idFatura); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// atualizarStatusFatura godoc
func atualizarStatusFatura(c *gin.Context) {
	var req faturas.ReqAtualizarStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	if err := faturas.AtualizarStatusFatura(&req, idFatura); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
