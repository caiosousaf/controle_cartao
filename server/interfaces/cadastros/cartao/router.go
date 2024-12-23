package cartao

import (
	"controle_cartao/interfaces/cadastros/faturas"
	"github.com/gin-gonic/gin"
)

// Router é um router para as rotas de cartões que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCartoes)
	r.POST("", cadastrarCartao)
}

// RouterWithID é um router para as rotas de cartões que utilizam ID
func RouterWithID(r *gin.RouterGroup) {
	r.GET("", buscarCartao)
	r.PUT("", atualizarCartao)
	r.DELETE("/remover", removerCartao)
	r.PUT("/reativar", reativarCartao)
	faturas.RouterCard(r.Group("faturas"))
	faturas.RouterWithCardID(r.Group("fatura/:fatura_id"))
}
