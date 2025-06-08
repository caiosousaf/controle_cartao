package usuarios

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.POST("", cadastrarUsuario)
	r.PUT("alterar/senha", atualizarSenhaUsuario)
}

func RouterWithID(r *gin.RouterGroup) {
	r.GET("", buscarUsuario)
}

func RouterLogin(r *gin.RouterGroup) {
	r.POST("login", login)
}
