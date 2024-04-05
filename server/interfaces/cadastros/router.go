package cadastros

import (
	"controle_cartao/interfaces/cadastros/cartao"
	"controle_cartao/interfaces/cadastros/categorias"
	"controle_cartao/interfaces/cadastros/compras"
	"controle_cartao/interfaces/cadastros/faturas"
	"github.com/gin-gonic/gin"
)

// Router é um router para gerenciamento das rotas de cadastros
func Router(r *gin.RouterGroup) {
	cartao.Router(r.Group("cartoes"))
	cartao.RouterWithID(r.Group("cartao/:cartao_id"))
	faturas.RouterWithID(r.Group("fatura/:fatura_id"))
	compras.Router(r.Group("compras"))
	categorias.Router(r.Group("categorias"))
}
