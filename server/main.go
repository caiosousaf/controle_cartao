package main

import (
	"controle_cartao/interfaces/cadastros"
	"controle_cartao/interfaces/cadastros/usuarios"
	"controle_cartao/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.Use(CORSMiddleware())

	RouterCadastro := r.Group("cadastros")
	RouterCadastro.Use(middlewares.AuthMiddleware())
	cadastros.Router(RouterCadastro)

	usuarios.RouterLogin(r.Group("usuarios"))

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ✅ Permite qualquer origem
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		// ✅ Trata requisições OPTIONS (pré-flight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
