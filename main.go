package main

import (
	"controle_cartao/interfaces/cadastros"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	cadastros.Router(r.Group("cadastros"))

	err := r.Run("localhost:8080")
	if err != nil {
		return
	}
}
