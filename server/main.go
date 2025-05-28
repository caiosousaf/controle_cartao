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
	allowedOrigins := map[string]bool{
		"https://kaleidoscopic-sfogliatella-b7f6c8.netlify.app":true,
		"https://marvelous-haupia-7c28da.netlify.app": true,
		"http://localhost:3000":                       true,
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
