package categorias

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de categorias que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCategorias)
	r.POST("", cadastrarCategoria)
}

// RouterWithID é um router para as rotas de categorias que utilizam ID
func RouterWithID(r *gin.RouterGroup) {
	r.DELETE("remover", removerCategoria)
	r.PUT("reativar", reativarCategoria)
	r.PUT("", atualizarCategoria)
	r.GET("", buscarCategoria)
}
