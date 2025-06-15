package compras

import (
	"controle_cartao/interfaces/cadastros/compras/recorrente"
	"github.com/gin-gonic/gin"
)

// Router é um router para as rotas de compras que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCompras)
	r.GET("total", obterTotalComprasValor)
	r.GET("pdf", pdfComprasFaturaCartao)
	recorrente.Router(r.Group("recorrente"))
	recorrente.RouterWithID(r.Group("recorrente/:recorrente_id"))
}

// RouterWithID é um router para as rotas de compras que utilizam id
func RouterWithID(r *gin.RouterGroup) {
	r.PUT(":remover_todas_parcelas", atualizarCompras)
	r.DELETE(":remover_todas_parcelas", removerCompra)
}

// RouterInvoice é um router para as rotas de compras que não utilizam ID de compra mas utiliza ID da fatura
func RouterInvoice(r *gin.RouterGroup) {
	r.POST("", cadastrarCompra)
}
