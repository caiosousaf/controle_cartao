package recorrente

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de compras recorrentes que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarComprasRecorrentes)
}
