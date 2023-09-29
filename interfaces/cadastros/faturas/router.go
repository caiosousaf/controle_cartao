package faturas

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de faturas que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarFaturasCartao)
	r.POST("", cadastrarFatura)
}

// RouterWithID é um router para as rotas de faturas que utilizam ID
func RouterWithID(r *gin.RouterGroup) {
	r.GET("", buscarFaturaCartao)
	r.PUT("", atualizarFatura)
}
