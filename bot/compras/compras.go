package compras

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

// ListarComprasFatura é responsável por realizar a requisição de listagem para compras
func ListarComprasFatura(idFatura *string) (res ResComprasPag) {
	resp, err := http.Get(BaseURLCompras + "?fatura_id=" + *idFatura)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	// Imprime a resposta da API
	fmt.Println("Resposta da API:", string(body))

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}

// ObterComprasPdf é responsável por realizar a requisição que obtém o pdf com as compras
func ObterComprasPdf(idFatura *uuid.UUID, idCartao *uuid.UUID) []byte {
	var (
		resp *http.Response
		err  error
	)

	if idFatura != nil && idCartao != nil {
		resp, err = http.Get(BaseURLComprasPdf + "?fatura_id=" + idFatura.String() + "&cartao_id=" + idCartao.String())
		if err != nil {
			fmt.Println("Erro ao fazer a requisição:", err)
			return nil
		}
		defer resp.Body.Close()
	} else {
		resp, err = http.Get(BaseURLComprasPdf + "?cartao_id=" + idCartao.String())
		if err != nil {
			fmt.Println("Erro ao fazer a requisição:", err)
			return nil
		}
		defer resp.Body.Close()
	}

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return nil
	}

	// Imprime a resposta da API
	fmt.Println("Resposta da API:", string(body))

	return body
}
