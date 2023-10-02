package compras

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de compras que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCompras)
}

// RouterInvoice é um router para as rotas de compras que não utilizam ID de compra mas utiliza ID da fatura
func RouterInvoice(r *gin.RouterGroup) {
	r.POST("", cadastrarCompra)
}
