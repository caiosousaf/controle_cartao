package recorrente

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de compras recorrentes que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarComprasRecorrentes)
	r.POST("", cadastrarComprasRecorrentes)
	r.GET("previsao", obterPrevisaoGastos)
	r.POST("cadastro", cadastrarNovaCompraRecorrente)
}

// RouterWithID é um router para as rotas de compras recorrente que utilizam ID
func RouterWithID(r *gin.RouterGroup) {
	r.PUT("desativar", desativarCompraRecorrente)
	r.PUT("reativar", reativarCompraRecorrente)
	r.DELETE("remover", removerCompraRecorrente)
	r.PUT("", atualizarCompraRecorrente)
}
