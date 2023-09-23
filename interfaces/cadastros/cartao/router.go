package cartao

import "github.com/gin-gonic/gin"

// Router é um router para as rotas de cartões que não utilizam ID
func Router(r *gin.RouterGroup) {
	r.GET("", listarCartoes)
}
