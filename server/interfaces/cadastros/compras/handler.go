package compras

import (
	"controle_cartao/application/cadastros/compras"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// cadastrarCompra godoc
func cadastrarCompra(c *gin.Context) {
	var req compras.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	idFatura, err := utils.GetUUIDFromParam(c, "fatura_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	idCompra, err := compras.CadastrarCompra(&req, idFatura, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, idCompra)
}

// listarCompras godoc
func listarCompras(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := compras.ListarCompras(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}

// pdfComprasFaturaCartao godoc
func pdfComprasFaturaCartao(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	pdf, err := compras.PdfComprasFaturaCartao(&params, usuarioID)
	if err != nil {
		fontPaths := []string{
			"/app/font/CaviarDreams.ttf",
			"/app/font/CaviarDreams_Bold.ttf",
			"/app/font/CaviarDreams.ttf",
			"/app/font/CaviarDreams_Bold.ttf",
			"app/font/CaviarDreams.ttf",
			"app/font/CaviarDreams_Bold.ttf",
			"/server/font/CaviarDreams.ttf",
			"/server/font/CaviarDreams_Bold.ttf",
			"server/font/CaviarDreams.ttf",
			"server/font/CaviarDreams_Bold.ttf",
			"font/CaviarDreams.ttf",
			"font/CaviarDreams_Bold.ttf",
			"/font/CaviarDreams.ttf",
			"/font/CaviarDreams_Bold.ttf",
		}

		for _, path := range fontPaths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				log.Printf("Font file not found: %s", path)
			}
		}
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	// Define o cabeçalho Content-Disposition para indicar que o arquivo é um anexo
	c.Header("Content-Disposition", "attachment; filename=compras.pdf")

	// Define o tipo de conteúdo como PDF
	c.Header("Content-Type", "application/pdf")

	// Gera o PDF e escreve no contexto de resposta
	log.Println(pdf.Output(c.Writer), pdf)
	if err = pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, "Erro ao gerar PDF")
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, "PDF gerado com sucesso")
}

// obterTotalComprasValor godoc
func obterTotalComprasValor(c *gin.Context) {
	params, err := utils.ParseParams(c)
	if err != nil {
		return
	}

	usuarioID := middlewares.AuthUsuario(c)

	res, err := compras.ObterTotalComprasValor(&params, usuarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error(err)})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res)
}
