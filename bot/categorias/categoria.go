package categorias

import (
	"bot_controle_cartao/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListarCategorias é responsável por realizar a requisição para listar as categorias
func ListarCategorias(url string, chatID int64, userTokens map[int64]string) (categorias ResCategoriasPag, err error) {
	var ambiente = utils.ValidarAmbiente()

	token, ok := userTokens[chatID]
	if !ok {
		return categorias, fmt.Errorf("usuário não está autenticado")
	}

	req, err := http.NewRequest(http.MethodGet, ambiente+url, nil)
	if err != nil {
		return categorias, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return categorias, fmt.Errorf("Realize login!")
	} else if resp.StatusCode != http.StatusOK {
		return categorias, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	if err = json.Unmarshal(body, &categorias); err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	return
}
