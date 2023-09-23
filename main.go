package main

import (
	"controle_cartao/interfaces/cadastros"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cadastros.Router(r.Group("cadastros"))

	r.Run("localhost:8080")
}
