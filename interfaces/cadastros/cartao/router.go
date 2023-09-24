package cartao

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de cartões que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCartoes)
	r.POST("", cadastrarCartao)
}

// RouterWithID é um router para as rotas de cartões que utilizam ID
func RouterWithID(r *gin.RouterGroup) {
	r.GET(":cartao_id", buscarCartao)
	r.PUT(":cartao_id", atualizarCartao)
}
