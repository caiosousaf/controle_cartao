package utils

import "time"

func ObterNumeroDoMes(dataStr string) (int, error) {
	// Parse a string no formato "YYYY-MM-DD" para uma data
	data, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		return 0, err // Retorna 0 e o erro se houver um problema na conversão
	}

	// Extraia o número do mês da data (1 a 12)
	mes := data.Month()

	return int(mes), nil
}

// NumeroParaNomeMes recebe um inteiro e retorna o mês referente aquele inteiro
func NumeroParaNomeMes(numero int) string {
	meses := map[int]string{
		1:  "Janeiro",
		2:  "Fevereiro",
		3:  "Março",
		4:  "Abril",
		5:  "Maio",
		6:  "Junho",
		7:  "Julho",
		8:  "Agosto",
		9:  "Setembro",
		10: "Outubro",
		11: "Novembro",
		12: "Dezembro",
	}

	nomeMes, ok := meses[numero]
	if ok {
		return nomeMes
	}
	return "Mês inválido"
}
