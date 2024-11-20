package usuarios

import "github.com/gin-gonic/gin"

func Router(r *gin.RouterGroup) {
	r.POST("", cadastrarUsuario)
}

func RouterLogin(r *gin.RouterGroup) {
	r.POST("login", login)
}
