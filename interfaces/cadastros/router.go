package cadastros

import (
	"controle_cartao/interfaces/cadastros/cartao"
	"github.com/gin-gonic/gin"
)

// Router é um router para gerenciamento das rotas de cadastros
func Router(r *gin.RouterGroup) {
	cartao.Router(r.Group("cartoes"))
}
