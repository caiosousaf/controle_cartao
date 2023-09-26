package utils

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
