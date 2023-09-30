package compras

import "github.com/gin-gonic/gin"

// RouterInvoice é um router para as rotas de compras que não utilizam ID de compra mas utiliza ID da fatura
func RouterInvoice(r *gin.RouterGroup) {
	r.POST("", cadastrarCompra)
}
