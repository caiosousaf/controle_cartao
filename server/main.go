package main

import (
	"controle_cartao/interfaces/cadastros"
	"controle_cartao/interfaces/cadastros/usuarios"
	"controle_cartao/middlewares"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
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
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := make(map[string]bool)

	if allowedOriginsStr != "" {
		origins := strings.Split(allowedOriginsStr, ",")
		for _, origin := range origins {
			allowedOrigins[strings.TrimSpace(origin)] = true
		}
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
