package utils

import "math"

// Min retorna o valor mínimo de um slice
func Min(valores ...int) int {
	if len(valores) == 0 {
		return math.MinInt32
	}

	min := valores[0]
	for _, v := range valores {
		if v < min {
			min = v
		}
	}

	return min
}

// ArredondarParaDuasCasasDecimais torna um float64 com várias casas decimais em apenas um float com 2 casas decimais
func ArredondarParaDuasCasasDecimais(valor *float64) {
	*valor = math.Round(*valor*100) / 100
}
