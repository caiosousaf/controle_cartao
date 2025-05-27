package main

import (
	"controle_cartao/interfaces/cadastros"
	"controle_cartao/interfaces/cadastros/usuarios"
	"controle_cartao/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	RouterCadastro := r.Group("cadastros")
	RouterCadastro.Use(middlewares.AuthMiddleware())
	cadastros.Router(RouterCadastro)

	usuarios.RouterLogin(r.Group("usuarios"))

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
