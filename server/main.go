package main

import (
	"controle_cartao/interfaces/cadastros"
	"controle_cartao/interfaces/cadastros/usuarios"
	"controle_cartao/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	RouterCadastro := r.Group("cadastros")
	RouterCadastro.Use(middlewares.AuthMiddleware())
	cadastros.Router(RouterCadastro)

	usuarios.RouterLogin(r.Group("usuarios"))

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
