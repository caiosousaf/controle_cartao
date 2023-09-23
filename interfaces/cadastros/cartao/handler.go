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
