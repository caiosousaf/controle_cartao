package faturas

import (
	"controle_cartao/interfaces/cadastros/compras"
	"github.com/gin-gonic/gin"
)

// RouterCard é um router para as rotas de faturas que não utilizam ID da fatura mas utilizam id do cartão
func RouterCard(r *gin.RouterGroup) {
	r.GET("", listarFaturasCartao)
	r.POST("", cadastrarFatura)
}

// RouterWithCardID é um router para as rotas de faturas que utilizam ID da fatura e do cartao
func RouterWithCardID(r *gin.RouterGroup) {
	r.PUT("", atualizarFatura)
}

// RouterWithID é um router para as rotas de faturas que utilizam ID da fatura
func RouterWithID(r *gin.RouterGroup) {
	r.PUT("status", atualizarStatusFatura)
	r.GET("", buscarFatura)
	compras.RouterInvoice(r.Group("compras"))
}
